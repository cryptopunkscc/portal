package project

import (
	"encoding/json"
	"github.com/leaanthony/debme"
	"github.com/pkg/errors"
	gofs "io/fs"
)

func init() {
	if err := loadTemplateCache(); err != nil {
		panic(err)
	}
}

// Loads the template cache
func loadTemplateCache() error {
	templatesFS, err := debme.FS(templatesFs, "tmpl")
	if err != nil {
		return errors.Wrap(err, "cannot open template fs")
	}

	// Get directories
	files, err := templatesFS.ReadDir(".")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			templateFS, err := templatesFS.FS(file.Name())
			if err != nil {
				return err
			}
			template, err := parseTemplate(templateFS)
			if err != nil {
				// Cannot parse this template, continue
				continue
			}
			templateCache = append(templateCache, template)
		}
	}

	return nil
}

func parseTemplate(template gofs.FS) (Template, error) {
	var result Template
	data, err := gofs.ReadFile(template, "template.json")
	if err != nil {
		return result, errors.Wrap(err, "Error parsing template")
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	result.FS = template
	return result, nil
}
