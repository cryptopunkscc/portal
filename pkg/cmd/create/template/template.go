package template

import (
	gofs "io/fs"
)

// Template holds data relating to a template
// including the metadata stored in template.json
type Template struct {
	// Template details
	Name        string `json:"name"`
	ShortName   string `json:"shortname"`
	Author      string `json:"author"`
	Description string `json:"description"`
	HelpURL     string `json:"helpurl"`

	// Other data
	FS gofs.FS `json:"-"`
}

// Data will be embedded into the tmpl files during the installation
type Data struct {
	ProjectName        string
	PackageName        string
	AuthorName         string
	AuthorEmail        string
	AuthorNameAndEmail string
	OutputFile         string
}
