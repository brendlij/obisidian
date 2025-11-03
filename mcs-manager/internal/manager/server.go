package manager

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/log"

	"obsidian/internal/query"
	"obsidian/internal/server"
	"obsidian/internal/util"
	"obsidian/pkg/events"
)

type ServerType = server.ServerType

const (
	TypeVanilla = server.TypeVanilla
	TypePaper   = server.TypePaper
)

type ServerConfig = server.ServerConfig

type ServerState string

const (
	StateStopped  ServerState = "stopped"
	StateRunning  ServerState = "running"
	StateStarting ServerState = "starting"
	StateCrashed  ServerState = "crashed"
)

type ServerInfo struct {
	Config      server.ServerConfig `json:"config"`
	State       ServerState         `json:"state"`
	PID         int                 `json:"pid"`
	UptimeSec   int64               `json:"uptimeSec"`
	LastExitErr string              `json:"lastExitErr"`
	Players     *PlayerInfo         `json:"players,omitempty"`
}

type PlayerInfo struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type Server struct {
	cfg     ServerConfig
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	logf    *os.File
	state   atomic.Value
	startAt time.Time
	lastErr string
}

func (s *Server) State() ServerState { return s.state.Load().(ServerState) }

func (s *Server) Info() ServerInfo {
	up := int64(0)
	if !s.startAt.IsZero() {
		up = int64(time.Since(s.startAt).Seconds())
	}
	pid := 0
	if s.cmd != nil && s.cmd.Process != nil {
		pid = s.cmd.Process.Pid
	}

	// Try to read player info if server is running
	var players *PlayerInfo
	if s.State() == StateRunning {
		// Try to ping the server directly on its port
		if status, err := query.PingServer("localhost", s.cfg.Port, 2*time.Second); err == nil {
			players = &PlayerInfo{Current: status.Players.Online, Max: status.Players.Max}
		} else {
			// Fallback: Try to read from logs if ping fails
			logPath := filepath.Join(s.cfg.Path, "mcs.log")
			if current, max, err := util.ReadPlayersFromLog(logPath); err == nil && (current > 0 || max > 0) {
				// If max is not set, try to read from config
				if max == 0 {
					if configMax, cfgErr := util.ReadMaxPlayersFromConfig(filepath.Join(s.cfg.Path, "server.properties")); cfgErr == nil {
						max = configMax
					}
				}
				players = &PlayerInfo{Current: current, Max: max}
			} else {
				// Last resort: read max from config only
				if max, cfgErr := util.ReadMaxPlayersFromConfig(filepath.Join(s.cfg.Path, "server.properties")); cfgErr == nil {
					players = &PlayerInfo{Current: 0, Max: max}
				}
			}
		}
	}

	return ServerInfo{Config: s.cfg, State: s.State(), PID: pid, UptimeSec: up, LastExitErr: s.lastErr, Players: players}
}

func (s *Server) Start(bus *events.Bus) error {
	if s.State() == StateRunning {
		log.Debug("server already running", "id", s.cfg.ID, "name", s.cfg.Name)
		return nil
	}
	log.Info("starting server", "id", s.cfg.ID, "name", s.cfg.Name, "port", s.cfg.Port)
	s.state.Store(StateStarting)
	jar := filepath.Join(s.cfg.Path, "server.jar")
	java := "java"
	if runtime.GOOS == "windows" {
		java = "java.exe"
	}
	args := []string{"-Xmx" + strconv.Itoa(s.cfg.MemoryMB) + "M", "-jar", jar, "nogui"}
	cmd := exec.Command(java, args...)
	cmd.Dir = s.cfg.Path
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdin, _ := cmd.StdinPipe()
	logFile, _ := os.OpenFile(filepath.Join(s.cfg.Path, "mcs.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	s.stdin, s.cmd, s.logf = stdin, cmd, logFile
	if err := cmd.Start(); err != nil {
		log.Error("failed to start server", "id", s.cfg.ID, "err", err)
		s.state.Store(StateCrashed)
		return err
	}
	log.Debug("server process started", "id", s.cfg.ID, "pid", cmd.Process.Pid)
	s.state.Store(StateRunning)
	s.startAt = time.Now()
	log.Info("server started successfully", "id", s.cfg.ID, "name", s.cfg.Name, "pid", cmd.Process.Pid)
	bus.Publish(events.Event{Type: "server.started", ServerID: s.cfg.ID})

	go s.pipe(bus, stdout, "stdout")
	go s.pipe(bus, stderr, "stderr")
	go func() {
		err := cmd.Wait()
		_ = logFile.Close()
		if err != nil {
			s.lastErr = err.Error()
			s.state.Store(StateCrashed)
			log.Error("server crashed", "id", s.cfg.ID, "err", err)
		} else {
			s.state.Store(StateStopped)
			log.Info("server stopped", "id", s.cfg.ID)
		}
		bus.Publish(events.Event{Type: "server.exited", ServerID: s.cfg.ID})
	}()
	return nil
}

func (s *Server) pipe(bus *events.Bus, r io.Reader, stream string) {
	scanner := bufio.NewScanner(r)
	// Use smaller buffer for more responsive logging (default is 65536)
	scanner.Buffer(make([]byte, 4096), 4096)
	for scanner.Scan() {
		line := scanner.Text()
		if s.logf != nil {
			_, _ = s.logf.WriteString(line + "\n")
		}
		bus.Publish(events.Event{Type: "server.log", ServerID: s.cfg.ID, Data: map[string]any{"stream": stream, "line": line}})
	}
}

func (s *Server) SendCommand(cmd string) error {
	if s.stdin == nil {
		log.Warn("server not running, cannot send command", "id", s.cfg.ID, "cmd", cmd)
		return errors.New("not running")
	}
	log.Debug("sending command to server", "id", s.cfg.ID, "cmd", cmd)
	_, err := io.WriteString(s.stdin, cmd+"\n")
	return err
}

func (s *Server) Stop(bus *events.Bus) {
	if s.State() != StateRunning {
		log.Debug("server not running, cannot stop", "id", s.cfg.ID)
		return
	}
	log.Info("stopping server", "id", s.cfg.ID, "name", s.cfg.Name)
	if s.stdin != nil {
		_, _ = io.WriteString(s.stdin, "stop\n")
	}
}

func (s *Server) Restart(bus *events.Bus) error {
	log.Info("restarting server", "id", s.cfg.ID, "name", s.cfg.Name)
	
	// If not running, just start it
	if s.State() != StateRunning {
		return s.Start(bus)
	}
	
	s.Stop(bus)
	
	// Wait for server to actually stop by listening to exit event
	sub := bus.Subscribe()
	defer bus.Unsubscribe(sub)
	
	// Wait with timeout (max 30 seconds)
	timeout := time.After(30 * time.Second)
	for {
		select {
		case e := <-sub.Ch:
			if e.Type == "server.exited" && e.ServerID == s.cfg.ID {
				log.Debug("server stopped successfully, starting again", "id", s.cfg.ID)
				return s.Start(bus)
			}
		case <-timeout:
			log.Warn("timeout waiting for server to stop, forcing start anyway", "id", s.cfg.ID)
			return s.Start(bus)
		}
	}
}
