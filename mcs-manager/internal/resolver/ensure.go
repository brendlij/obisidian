package resolver

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"mcs-manager/internal/server"
)

func EnsureJar(cfg server.ServerConfig, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		return nil
	}
	url := cfg.JarURL
	if url == "" {
		var err error
		switch cfg.Type {
		case server.TypeVanilla:
			url, err = resolveVanilla(cfg.Version)
		case server.TypePaper:
			url, err = resolvePaper(cfg.Version)
		default:
			return fmt.Errorf("resolver: type %s not supported without jarUrl", cfg.Type)
		}
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	return downloadTo(url, io.Writer(f))
}
