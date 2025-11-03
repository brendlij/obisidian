package server

type ServerType string

const (
	TypeVanilla ServerType = "vanilla"
	TypePaper   ServerType = "paper"
	TypeFabric  ServerType = "fabric"
)

type ServerConfig struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     ServerType `json:"type"`
	Version  string     `json:"version"`
	Port     int        `json:"port"`
	MemoryMB int        `json:"memoryMb"`
	Path     string     `json:"path"`
	Eula     bool       `json:"eula"`
	JarURL   string     `json:"jarUrl"`
}
