package create

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target/template"
	"github.com/leaanthony/gosod"
	"io/fs"
)

func Install(dir string, template fs.FS, args template.Args) (err error) {
	if err = extractTemplate(dir, template, args); err != nil {
		return fmt.Errorf("extract template: %w", err)
	}
	if err = extractCommons(dir, args); err != nil {
		return fmt.Errorf("extract commons: %w", err)
	}
	return
}

func extractTemplate(dir string, files fs.FS, args template.Args) (err error) {
	installer := gosod.New(files)
	installer.IgnoreFile("template.json")
	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("template.Extract: %v", err)
	}
	return
}

func extractCommons(dir string, args template.Args) (err error) {
	installer := gosod.New(template.CommonsFs)
	if err = installer.Extract(dir, args); err != nil {
		return fmt.Errorf("template.Extract: %v", err)
	}
	return
}
