package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	"obsidian/internal/api"
	"obsidian/internal/manager"
	"obsidian/internal/store"
	"obsidian/pkg/events"
)

type BootConfig struct {
	Root string `json:"root"`
	Bind string `json:"bind"`
}

func main() {
	cfg := BootConfig{Root: defaultRoot(), Bind: ":8484"}
	if env := os.Getenv("MCS_CONFIG"); env != "" {
		log.Info("loading config from env", "path", env)
		b, err := os.ReadFile(env)
		if err == nil {
			_ = json.Unmarshal(b, &cfg)
			log.Info("config loaded", "root", cfg.Root, "bind", cfg.Bind)
		} else {
			log.Warn("failed to load config file", "path", env, "err", err)
		}
	}

	log.Info("initializing manager", "root", cfg.Root)
	bus := events.NewBus()
	st := store.NewJSON(cfg.Root)
	mgr, err := manager.New(cfg.Root, bus, st)
	if err != nil {
		log.Fatal("failed to initialize manager", "err", err)
	}

	log.Info("starting HTTP API server", "bind", cfg.Bind)
	apiSrv := api.NewHTTP(cfg.Bind, mgr, bus)
	log.Fatal(apiSrv.ListenAndServe())
}

func defaultRoot() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "mcs-servers")
}
