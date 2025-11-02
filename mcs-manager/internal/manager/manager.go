package manager

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"obsidian/internal/resolver"
	"obsidian/internal/server"
	"obsidian/internal/util"
	"obsidian/pkg/events"
)

type Manager struct {
	root  string
	mu    sync.RWMutex
	items map[string]*Server
	bus   *events.Bus
	store Store
}

type Store interface {
	LoadAll() ([]ServerConfig, error)
	SaveAll([]ServerConfig) error
}

func New(root string, bus *events.Bus, st Store) (*Manager, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, err
	}
	m := &Manager{root: root, items: map[string]*Server{}, bus: bus, store: st}

	// Load persisted servers
	if servers, err := st.LoadAll(); err == nil {
		for _, cfg := range servers {
			s := &Server{cfg: cfg}
			s.state.Store(StateStopped)
			m.items[cfg.ID] = s
		}
	}

	return m, nil
}

func (m *Manager) List() []ServerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]ServerInfo, 0, len(m.items))
	for _, s := range m.items {
		out = append(out, s.Info())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Config.Name < out[j].Config.Name })
	return out
}

func (m *Manager) Get(id string) (*Server, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.items[id]
	return s, ok
}

func (m *Manager) Create(cfg server.ServerConfig) (*Server, error) {
	if cfg.ID == "" {
		cfg.ID = util.RandID()
	}
	if cfg.Name == "" {
		cfg.Name = cfg.ID
	}
	if cfg.Path == "" {
		cfg.Path = filepath.Join(m.root, cfg.ID)
	}
	if cfg.MemoryMB == 0 {
		cfg.MemoryMB = 2048
	}
	if cfg.Port == 0 {
		p, _ := util.PickFreePort()
		cfg.Port = p
	}
	if cfg.Version == "" {
		cfg.Version = "latest"
	}
	if cfg.Type == "" {
		cfg.Type = TypeVanilla
	}
	if err := os.MkdirAll(cfg.Path, 0o755); err != nil {
		return nil, err
	}
	if cfg.Eula {
		_ = os.WriteFile(filepath.Join(cfg.Path, "eula.txt"), []byte("eula=true\n"), 0o644)
	}
	// Jar download
	jarPath := filepath.Join(cfg.Path, "server.jar")
	if err := resolver.EnsureJar(cfg, jarPath); err != nil {
		return nil, err
	}

	s := &Server{cfg: cfg}
	s.state.Store(StateStopped)
	m.mu.Lock()
	m.items[cfg.ID] = s
	m.mu.Unlock()
	m.bus.Publish(events.Event{Type: "server.created", ServerID: cfg.ID, Data: cfg})
	_ = m.persist()
	return s, nil
}

func (m *Manager) Delete(id string) error {
	s, ok := m.Get(id)
	if !ok {
		return os.ErrNotExist
	}
	if s.State() == StateRunning {
		return errors.New("server running")
	}
	m.mu.Lock()
	delete(m.items, id)
	m.mu.Unlock()
	_ = os.RemoveAll(s.cfg.Path)
	_ = m.persist()
	m.bus.Publish(events.Event{Type: "server.deleted", ServerID: id})
	return nil
}

func (m *Manager) persist() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := []ServerConfig{}
	for _, s := range m.items {
		arr = append(arr, s.cfg)
	}
	return m.store.SaveAll(arr)
}
