package template

// Args will be embedded into the tmpl files during the installation
type Args struct {
	ProjectName string
	PackageName string
	AuthorName  string
	AuthorEmail string
	Description string
	Url         string
}
