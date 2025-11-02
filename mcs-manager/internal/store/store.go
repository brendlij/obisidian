package store

import "mcs-manager/internal/manager"

type Store interface {
	LoadAll() ([]manager.ServerConfig, error)
	SaveAll([]manager.ServerConfig) error
}
