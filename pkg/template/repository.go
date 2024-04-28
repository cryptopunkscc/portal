package template

import (
	"embed"
	"fmt"
)

//go:embed all:tmpl
var templatesFs embed.FS

// Cache for the templates
// We use this because we need different views of the same data
var templateCache []Template

// Templates returns the list of available templates
func Templates() ([]Template, error) {
	return templateCache, nil
}

// getTemplateByShortname returns the template with the given short name
func getTemplateByShortname(shortname string) (result Template, err error) {
	for _, result = range templateCache {
		if result.ShortName == shortname {
			return
		}
	}
	return Template{}, fmt.Errorf("shortname '%s' is not a valid template shortname", shortname)
}
