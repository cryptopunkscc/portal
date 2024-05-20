package target

const PackageJsonFilename = "package.json"

type PackageJson struct {
	Module  string `json:"module,omitempty"`
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func (p PackageJson) IsPortalLib() bool {
	return p.Module == "portal"
}

func (p PackageJson) CanBuild() bool {
	return p.Scripts.Build != ""
}
