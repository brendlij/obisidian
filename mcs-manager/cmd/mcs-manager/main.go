package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"mcs-manager/internal/api"
	"mcs-manager/internal/manager"
	"mcs-manager/internal/store"
	"mcs-manager/pkg/events"
)

type BootConfig struct {
	Root string `json:"root"`
	Bind string `json:"bind"`
}

func main() {
	cfg := BootConfig{Root: defaultRoot(), Bind: ":8484"}
	if env := os.Getenv("MCS_CONFIG"); env != "" {
		b, err := os.ReadFile(env)
		if err == nil {
			_ = json.Unmarshal(b, &cfg)
		}
	}

	bus := events.NewBus()
	st := store.NewJSON(cfg.Root)
	mgr, err := manager.New(cfg.Root, bus, st)
	if err != nil {
		log.Fatal(err)
	}

	apiSrv := api.NewHTTP(cfg.Bind, mgr, bus)
	log.Printf("[MCS] Manager running on %s", cfg.Bind)
	log.Fatal(apiSrv.ListenAndServe())
}

func defaultRoot() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "mcs-servers")
}
