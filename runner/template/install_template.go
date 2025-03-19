package template

import (
	"fmt"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/leaanthony/gosod"
	"io/fs"
)

func installTemplate(dir string, template fs.FS, args Args) (err error) {
	if err = extractTemplate(dir, template, args); err == nil {
		if err = extractCommons(dir, args); err == nil {
			return
		}
	}
	return fmt.Errorf("cannot install template: %w", err)
}

func extractTemplate(dir string, files fs.FS, args Args) (err error) {
	installer := gosod.New(files)
	installer.IgnoreFile("template.json")
	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("cannot extract template: %v", err)
	}
	return
}

func extractCommons(dir string, args Args) (err error) {
	installer := gosod.New(template.CommonsFs)
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("cannot extract template commons: %v", err)
	}
	return
}
