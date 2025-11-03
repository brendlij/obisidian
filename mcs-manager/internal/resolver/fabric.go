package resolver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const fabricMetaAPI = "https://meta.fabricmc.net/v2/versions"

type FabricGameVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

type FabricLoaderVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

type FabricInstallerVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// GetFabricVersions returns available Fabric-supported Minecraft versions
func GetFabricVersions() ([]string, error) {
	resp, err := http.Get(fabricMetaAPI + "/game")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fabric game versions: %w", err)
	}
	defer resp.Body.Close()

	var versions []FabricGameVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode fabric versions: %w", err)
	}

	// Filter only stable versions and convert to strings
	var result []string
	for _, v := range versions {
		if v.Stable {
			result = append(result, v.Version)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no stable fabric versions available")
	}

	return result, nil
}

// GetFabricLoaderVersions returns available Fabric loader versions
func GetFabricLoaderVersions() ([]string, error) {
	resp, err := http.Get(fabricMetaAPI + "/loader")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fabric loader versions: %w", err)
	}
	defer resp.Body.Close()

	var versions []FabricLoaderVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode fabric loader versions: %w", err)
	}

	// Get first stable version (latest)
	for _, v := range versions {
		if v.Stable {
			return []string{v.Version}, nil
		}
	}

	return nil, fmt.Errorf("no stable loader versions available")
}

// GetFabricInstallerVersion returns the latest stable installer version
func GetFabricInstallerVersion() (string, error) {
	resp, err := http.Get(fabricMetaAPI + "/installer")
	if err != nil {
		return "", fmt.Errorf("failed to fetch fabric installer versions: %w", err)
	}
	defer resp.Body.Close()

	var versions []FabricInstallerVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", fmt.Errorf("failed to decode fabric installer versions: %w", err)
	}

	// Get first stable version (latest)
	for _, v := range versions {
		if v.Stable {
			return v.Version, nil
		}
	}

	return "", fmt.Errorf("no stable installer versions available")
}

// ResolveFabric returns the download URL for a Fabric server JAR
func ResolveFabric(mcVersion string) (string, error) {
	// Get latest stable loader version
	loaderVersions, err := GetFabricLoaderVersions()
	if err != nil {
		return "", err
	}

	if len(loaderVersions) == 0 {
		return "", fmt.Errorf("no loader versions available")
	}

	loaderVersion := loaderVersions[0]

	// Get latest stable installer version
	installerVersion, err := GetFabricInstallerVersion()
	if err != nil {
		return "", err
	}

	// Construct the download URL
	downloadURL := fmt.Sprintf(
		"%s/loader/%s/%s/%s/server/jar",
		fabricMetaAPI,
		mcVersion,
		loaderVersion,
		installerVersion,
	)

	return downloadURL, nil
}

