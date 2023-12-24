package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type fwVer struct {
	From string `json:"from"`
}
type Manifest struct {
	Name            string  `json:"name"`
	FirmwareVersion []fwVer `json:"firmware_versions"`
	AppLog          string  `json:"app_log"`
	Version         string  `json:"version"`
	Build           string  `json:"build"`
}

func CreateManifest(name, logfile string) Manifest {
	return Manifest{
		Name:            name,
		FirmwareVersion: []fwVer{{"2.1.2"}},
		AppLog:          logfile,
		Version:         "0.0.1",
		Build:           "1",
	}
}

func GetManifest(dir string) (Manifest, error) {
	data, err := os.ReadFile(path.Join(dir, "manifest.json"))
	if err != nil {
		return Manifest{}, err
	}
	var manifest Manifest
	err = json.Unmarshal(data, &manifest)
	return manifest, nil
}

func (m Manifest) GetVersionString() string {
	var ver string
	if m.Version == "" {
		return "v0.0.0"
	} else {
		ver = "v"
		ver += m.Version
		if m.Build != "" {
			ver += fmt.Sprintf("-%v", m.Build)
		}
	}
	return ver
}
