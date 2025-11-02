package resolver

import (
	"fmt"
)

type paperProject struct {
	Versions []string `json:"versions"`
}

type paperBuilds struct {
	Builds []struct {
		Build int `json:"build"`
	} `json:"builds"`
}

type paperArtifact struct {
	Downloads struct {
		Application struct {
			Name string `json:"name"`
		} `json:"application"`
	} `json:"downloads"`
}

func resolvePaper(version string) (string, error) {
	ver := version
	if ver == "" || ver == "latest" {
		var meta paperProject
		if err := getJSON("https://api.papermc.io/v2/projects/paper", &meta); err != nil {
			return "", err
		}
		if len(meta.Versions) == 0 {
			return "", fmt.Errorf("paper: no versions")
		}
		ver = meta.Versions[len(meta.Versions)-1]
	}
	var builds paperBuilds
	if err := getJSON(fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds", ver), &builds); err != nil {
		return "", err
	}
	if len(builds.Builds) == 0 {
		return "", fmt.Errorf("paper: no builds for %s", ver)
	}
	build := builds.Builds[len(builds.Builds)-1].Build
	var art paperArtifact
	if err := getJSON(fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d", ver, build), &art); err != nil {
		return "", err
	}
	name := art.Downloads.Application.Name
	return fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/%s", ver, build, name), nil
}
