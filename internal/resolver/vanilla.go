package resolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type mojangManifest struct {
	Latest   struct{ Release, Snapshot string } `json:"latest"`
	Versions []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"versions"`
}

type mojangVersion struct {
	Downloads struct {
		Server struct {
			Url string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func resolveVanilla(version string) (string, error) {
	const manifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
	var man mojangManifest
	if err := getJSON(manifestURL, &man); err != nil {
		return "", err
	}
	ver := version
	if ver == "" || ver == "latest" || ver == "release" {
		ver = man.Latest.Release
	}
	var vURL string
	for _, v := range man.Versions {
		if v.ID == ver {
			vURL = v.Url
			break
		}
	}
	if vURL == "" {
		return "", fmt.Errorf("vanilla: version not found: %s", ver)
	}
	var vd mojangVersion
	if err := getJSON(vURL, &vd); err != nil {
		return "", err
	}
	if vd.Downloads.Server.Url == "" {
		return "", errors.New("vanilla: server jar not available")
	}
	return vd.Downloads.Server.Url, nil
}

func httpGet(url string) (*http.Response, error) {
	return http.Get(url)
}

func getJSON(url string, out any) error {
	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("http %d for %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func downloadTo(url string, w io.Writer) error {
	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download http %d", resp.StatusCode)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}
