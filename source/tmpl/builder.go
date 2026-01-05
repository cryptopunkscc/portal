package tmpl

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/leaanthony/gosod"
	"github.com/spf13/afero"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed templates
var templatesFS embed.FS

type CreateOpts struct {
	Name         string `cli:"name n"`
	Template     string `cli:"template t"`
	TemplateArgs TemplateArgs
}

type TemplateArgs map[string]string

func Create(opts CreateOpts, path string) (err error) {
	return Builder{CreateOpts: opts}.WriteRef(*source.OSRef(path))
}

type Builder struct {
	source.Ref
	fs fs.FS
	CreateOpts
	Manifest  app.Manifest
	PortalYml string
}

var _ source.Writer = &Builder{}

func (b Builder) WriteRef(ref source.Ref) (err error) {
	defer plog.TraceErr(&err)
	if err = b.prepare(); err != nil {
		return
	}
	if err = checkTargetDirNotExist(ref.Path); err != nil {
		return
	}
	if err = makeSourceDir(ref.Path); err != nil {
		return
	}
	if err = b.extractTemplate(ref.Path); err != nil {
		return
	}
	if err = b.writeManifest(ref.Path); err != nil {
		return
	}
	return
}

func (b *Builder) prepare() (err error) {
	if b.Template == "" {
		return fmt.Errorf("template name required")
	}
	if b.fs, err = fs.Sub(templatesFS, path.Join("templates", b.Template)); err != nil {
		return
	}
	b.Fs = &afero.FromIOFS{FS: b.fs}

	if _, err = b.Fs.Stat("dev"); err == nil {
		b.PortalYml = "dev.portal.yml"
	} else {
		b.PortalYml = "portal.yml"
	}
	if b.Name == "" {
		b.Name = b.Template
	}
	b.Manifest.Name = b.Name
	if b.Manifest.Title == "" {
		b.Manifest.Title = cases.Title(language.English).String(b.Manifest.Name)
	}
	if b.Manifest.Description == "" {
		b.Manifest.Description = "missing description"
	}
	if b.Manifest.Package == "" {
		b.Manifest.Package = strings.Join([]string{"my", "app", b.Manifest.Name}, ".")
	}
	return
}

func (b Builder) writeManifest(dir string) (err error) {
	defer plog.TraceErr(&err)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	portalYml := filepath.Join(dir, b.PortalYml)
	bytes, err := yaml.Marshal(b.Manifest)
	if err != nil {
		return
	}
	if err = os.WriteFile(portalYml, bytes, 0755); err != nil {
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

func (b Builder) extractTemplate(dir string) (err error) {
	defer plog.TraceErr(&err)
	installer := gosod.New(b.fs)
	installer.IgnoreFile("template.json")
	installer.IgnoreFile("dev")
	installer.RenameFiles(map[string]string{
		"gitignore.txt": ".gitignore",
	})
	if err = installer.Extract(dir, b.TemplateArgs); err != nil {
		return fmt.Errorf("cannot extract template: %v", err)
	}
	return
}
