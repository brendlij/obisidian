package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"mcs-manager/internal/manager"
)

type jsonStore struct {
	root string
	file string
}

func NewJSON(root string) Store {
	return &jsonStore{root: root, file: filepath.Join(root, "servers.json")}
}

func (s *jsonStore) LoadAll() ([]manager.ServerConfig, error) {
	b, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return []manager.ServerConfig{}, nil
		}
		return nil, err
	}
	var arr []manager.ServerConfig
	if err := json.Unmarshal(b, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func (s *jsonStore) SaveAll(arr []manager.ServerConfig) error {
	b, err := json.MarshalIndent(arr, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.file + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.file)
}
