package resolver

import "errors"

func resolveFabric(version string) (string, error) {
	// TODO: Implement Fabric API resolution.
	// Workaround: set cfg.JarURL beim Anlegen des Servers.
	return "", errors.New("fabric resolver TODO – provide jarUrl in config")
}
