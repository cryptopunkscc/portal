package project

import (
	"encoding/json"
	"io/fs"
)

type PackageJson struct {
	Module  string `json:"module,omitempty"`
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func LoadPackageJson(files fs.FS) (pkgJson PackageJson, err error) {
	file, err := fs.ReadFile(files, "package.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &pkgJson)
	return
}

func (p PackageJson) HasBuildScript() bool {
	return p.Scripts.Build != ""
}

func (p PackageJson) IsPortalLib() bool {
	return p.Module == "portal"
}
