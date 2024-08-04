package target

import "strings"

const ManifestFilename = "portal"

type Manifest struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Package     string `json:"package,omitempty" yaml:"package,omitempty"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	Exec        string `json:"exec,omitempty" yaml:"exec,omitempty"`
	Env         Env    `json:"env,omitempty" yaml:"env,omitempty"`
}

func (m Manifest) Match(id string) bool {
	return id == m.Name || strings.HasPrefix(id, m.Package)
}

type Env struct {
	Timeout int64 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}
