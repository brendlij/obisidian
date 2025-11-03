package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"obsidian/internal/manager"
	"obsidian/pkg/events"
)

type API struct {
	mgr *manager.Manager
	bus *events.Bus
}

func NewHTTP(bind string, mgr *manager.Manager, bus *events.Bus) *http.Server {
	api := &API{mgr: mgr, bus: bus}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/servers", api.handleServers)
	mux.HandleFunc("/servers/", api.handleServerByID)
	mux.HandleFunc("/events", api.handleSSE)
	return &http.Server{Addr: bind, Handler: withCORS(mux)}
}

func (a *API) handleServers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Debug("listing servers")
		writeJSON(w, a.mgr.List())
	case http.MethodPost:
		log.Info("creating new server from API request")
		var cfg manager.ServerConfig
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			log.Error("failed to decode server config", "err", err)
			http.Error(w, err.Error(), 400)
			return
		}
		s, err := a.mgr.Create(cfg)
		if err != nil {
			log.Error("failed to create server", "err", err)
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, s.Info())
	default:
		w.WriteHeader(405)
	}
}

func (a *API) handleServerByID(w http.ResponseWriter, r *http.Request) {
	tail := strings.TrimPrefix(r.URL.Path, "/servers/")
	tail = strings.TrimSuffix(tail, "/")
	parts := strings.Split(tail, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id := parts[0]

	s, ok := a.mgr.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	// /servers/{id}
	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, s.Info())
			return
		case http.MethodDelete:
			if err := a.mgr.Delete(id); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			w.WriteHeader(204)
			return
		default:
			w.WriteHeader(405)
			return
		}
	}

	// /servers/{id}/{action}
	action := parts[1]
	switch action {
	case "start":
		if r.Method != http.MethodPost {
			w.WriteHeader(405); return
		}
		log.Info("API request to start server", "id", id)
		if err := s.Start(a.bus); err != nil {
			http.Error(w, err.Error(), 400); return
		}
		writeJSON(w, s.Info())
	case "stop":
		if r.Method != http.MethodPost {
			w.WriteHeader(405); return
		}
		log.Info("API request to stop server", "id", id)
		s.Stop(a.bus)
		writeJSON(w, s.Info())
	case "restart":
		if r.Method != http.MethodPost {
			w.WriteHeader(405); return
		}
		log.Info("API request to restart server", "id", id)
		if err := s.Restart(a.bus); err != nil {
			http.Error(w, err.Error(), 400); return
		}
		writeJSON(w, s.Info())
	case "cmd":
		if r.Method != http.MethodPost {
			w.WriteHeader(405); return
		}
		var body struct{ Command string `json:"command"` }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400); return
		}
		log.Info("API request to send command", "id", id, "cmd", body.Command)
		if err := s.SendCommand(body.Command); err != nil {
			http.Error(w, err.Error(), 400); return
		}
		w.WriteHeader(204)
	case "logs":
		if r.Method != http.MethodGet { w.WriteHeader(405); return }
		log.Debug("fetching server logs", "id", id)
		base := s.Info().Config.Path
		candidates := []string{
			filepath.Join(base, "mcs.log"),
			filepath.Join(base, "mcs-manager.log"), // fallback für ältere Builds
		}
		var text string
		var err error
		for _, p := range candidates {
			if _, statErr := os.Stat(p); statErr == nil {
				text, err = tailLines(p, 200)
				break
			}
		}
		if err != nil {
			http.Error(w, err.Error(), 500); return
		}
		if text == "" {
			w.WriteHeader(204); return
		}
		w.Header().Set("Content-Type","text/plain")
		_, _ = w.Write([]byte(text))
		return
	default:
		http.NotFound(w, r)
	}
}

// kleines Tail (ohne extra util-Import)
func tailLines(path string, n int) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(b), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n"), nil
}


func (a *API) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	sub := a.bus.Subscribe()
	defer a.bus.Unsubscribe(sub)
	for {
		select {
		case <-r.Context().Done():
			return
		case ev := <-sub.Ch:
			b, _ := json.Marshal(ev)
			w.Write([]byte("event: " + ev.Type + "\n"))
			w.Write([]byte("data: " + string(b) + "\n\n"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}
