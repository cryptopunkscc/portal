package target

import "strings"

const PortalJsonFilename = "portal.json"

type Manifest struct {
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Package     string `json:"package,omitempty"`
	Version     string `json:"version,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Exec        string `json:"exec,omitempty"`
	Build       string `json:"build,omitempty"`
	Env         Env    `json:"env,omitempty"`
}

func (m Manifest) Match(id string) bool {
	return id == m.Name || strings.HasPrefix(id, m.Package)
}

type Env struct {
	Timeout int64 `json:"timeout,omitempty"`
}
