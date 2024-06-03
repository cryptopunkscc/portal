package target

const PortalJsonFilename = "portal.json"

type Manifest struct {
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Package     string `json:"package,omitempty"`
	Version     string `json:"version,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Exec        string `json:"exec,omitempty"`
}
