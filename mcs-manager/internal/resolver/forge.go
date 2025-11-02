package resolver

import "errors"

func resolveForge(version string) (string, error) {
	// TODO: Implement Forge installer resolution.
	// Workaround: set cfg.JarURL beim Anlegen des Servers.
	return "", errors.New("forge resolver TODO – provide jarUrl in config")
}
