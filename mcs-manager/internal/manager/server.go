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

	"obsidian/internal/server"
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
	return ServerInfo{Config: s.cfg, State: s.State(), PID: pid, UptimeSec: up, LastExitErr: s.lastErr}
}

func (s *Server) Start(bus *events.Bus) error {
	if s.State() == StateRunning {
		return nil
	}
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
		s.state.Store(StateCrashed)
		return err
	}
	s.state.Store(StateRunning)
	s.startAt = time.Now()
	bus.Publish(events.Event{Type: "server.started", ServerID: s.cfg.ID})

	go s.pipe(bus, stdout, "stdout")
	go s.pipe(bus, stderr, "stderr")
	go func() {
		err := cmd.Wait()
		_ = logFile.Close()
		if err != nil {
			s.lastErr = err.Error()
			s.state.Store(StateCrashed)
		} else {
			s.state.Store(StateStopped)
		}
		bus.Publish(events.Event{Type: "server.exited", ServerID: s.cfg.ID})
	}()
	return nil
}

func (s *Server) pipe(bus *events.Bus, r io.Reader, stream string) {
	scanner := bufio.NewScanner(r)
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
		return errors.New("not running")
	}
	_, err := io.WriteString(s.stdin, cmd+"\n")
	return err
}

func (s *Server) Stop(bus *events.Bus) {
	if s.State() != StateRunning {
		return
	}
	if s.stdin != nil {
		_, _ = io.WriteString(s.stdin, "stop\n")
	}
}
