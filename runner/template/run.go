package template

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/git"
	"github.com/cryptopunkscc/portal/resolve/template"
	"io/fs"
	"os"
	"path/filepath"
)

type runner struct {
	dir       string
	templates map[string]string
	args      template.Args
}

func Runner(dir string, templates map[string]string) target.Run[target.Template] {
	args := template.Args{}
	args.AuthorName = git.UserName(false)
	if args.AuthorName == "" {
		args.AuthorName = git.UserName(true)
	}
	args.AuthorEmail = git.UserEmail(false)
	if args.AuthorEmail == "" {
		args.AuthorEmail = git.UserEmail(true)
	}
	r := &runner{
		dir:       target.Abs(dir),
		templates: templates,
		args:      args,
	}
	return r.Run
}

func (r *runner) Run(_ context.Context, t target.Template, _ ...string) (err error) {
	name := r.templates[t.Name()]
	dir := filepath.Join(r.dir, name)

	if err = r.checkTargetDir(dir, name); err != nil {
		return fmt.Errorf("check target dir: %w", err)
	}
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	args := r.targetArgs(t)
	if err = Install(dir, t.Files(), args); err != nil {
		err = fmt.Errorf("cannot install template: %w", err)
	}
	if isDev(t) {
		if err = os.Rename(
			filepath.Join(dir, "portal.json"),
			filepath.Join(dir, "dev.portal.json"),
		); err != nil {
			return
		}
	}
	return
}

// check if already exist
func (r *runner) checkTargetDir(dir, name string) (err error) {
	_, err = os.Stat(dir)
	if err == nil {
		err = fmt.Errorf("cannot create project %s: %s already exists", name, dir)
	} else if os.IsNotExist(err) {
		err = nil
	} else {
		err = fmt.Errorf("cannot create project %s: %v", name, err)
	}
	return
}

func (r *runner) targetArgs(t target.Template) (args template.Args) {
	name := r.templates[t.Name()]
	args = r.args
	args.ProjectName = name
	args.PackageName = "new.portal." + name
	return
}

func isDev(t target.Template) bool {
	stat, err := fs.Stat(t.Files(), "dev")
	return err == nil && !stat.IsDir()
}
