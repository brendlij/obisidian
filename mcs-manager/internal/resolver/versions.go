package resolver

import (
	"fmt"
)

// GetPaperVersions fetches available Paper versions from the Paper API
func GetPaperVersions() ([]string, error) {
	var meta paperProject
	if err := getJSON("https://api.papermc.io/v2/projects/paper", &meta); err != nil {
		return nil, fmt.Errorf("failed to fetch Paper versions: %w", err)
	}
	if len(meta.Versions) == 0 {
		return nil, fmt.Errorf("no Paper versions available")
	}
	return meta.Versions, nil
}

// GetVanillaVersions fetches available Vanilla versions from the Mojang API
func GetVanillaVersions() ([]string, error) {
	const manifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	var man mojangManifest
	if err := getJSON(manifestURL, &man); err != nil {
		return nil, fmt.Errorf("failed to fetch Vanilla versions: %w", err)
	}

	versions := make([]string, 0, len(man.Versions))
	for _, v := range man.Versions {
		// Only include release versions
		if v.Type == "release" {
			versions = append(versions, v.ID)
		}
	}
	return versions, nil
}

