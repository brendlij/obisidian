package store

import "mcs/internal/manager"

type Store interface {
	LoadAll() ([]manager.ServerConfig, error)
	SaveAll([]manager.ServerConfig) error
}
