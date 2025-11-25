package source

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/leaanthony/gosod"
)

type CreateOpts struct {
	TemplateOpts
	TemplatesFS fs.FS
	Path        string
}

type TemplateOpts struct {
	Template     string `cli:"template t"`
	TemplateArgs TemplateArgs
}

func Create(opts CreateOpts) (err error) {
	if opts.Template == "" {
		return
	}
	defer plog.TraceErr(&err)
	templateFS, err := fs.Sub(opts.TemplatesFS, opts.Template)
	if err != nil {
		return
	}
	if opts.Path == "" {
		return fmt.Errorf("source path required")
	}
	if err = checkTargetDirNotExist(opts.Path); err != nil {
		return
	}
	if err = makeSourceDir(opts.Path); err != nil {
		return
	}
	if err = extractTemplate(opts.Path, templateFS, opts.TemplateArgs); err != nil {
		return
	}
	return
}

func checkTargetDirNotExist(path string) (err error) {
	defer plog.TraceErr(&err)
	_, err = os.Stat(path)
	switch {
	case err == nil:
		return fmt.Errorf("cannot create project %s already exists", path)
	case os.IsNotExist(err):
		return nil
	default:
		return fmt.Errorf("cannot create project %s: %v", path, err)
	}
}

func makeSourceDir(dir string) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot create project dir %s: %w", dir, err)
	}
	return
}

type TemplateArgs map[string]string

func extractTemplate(dir string, files fs.FS, args TemplateArgs) (err error) {
	defer plog.TraceErr(&err)
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
