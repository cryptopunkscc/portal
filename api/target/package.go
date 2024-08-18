package target

const PackageJsonFilename = "package"

type PackageJson struct {
	Portal  string `json:"portal,omitempty"`
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func (p PackageJson) IsPortalLib() bool {
	return p.Portal == "lib"
}

func (p PackageJson) CanBuild() bool {
	return p.Scripts.Build != ""
}
