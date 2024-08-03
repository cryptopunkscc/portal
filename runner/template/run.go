package template

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/git"
	"github.com/cryptopunkscc/portal/resolve/template"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"path/filepath"
)

func NewRun(dir string, templates map[string]string) target.Run[target.Template] {
	return NewRunner(dir, templates).Run
}

type Runner struct {
	dir       string
	templates map[string]string
	args      template.Args
}

func NewRunner(dir string, templates map[string]string) *Runner {
	args := template.Args{}
	args.AuthorName = git.UserName(false)
	if args.AuthorName == "" {
		args.AuthorName = git.UserName(true)
	}
	args.AuthorEmail = git.UserEmail(false)
	if args.AuthorEmail == "" {
		args.AuthorEmail = git.UserEmail(true)
	}
	return &Runner{
		dir:       target.Abs(dir),
		templates: templates,
		args:      args,
	}
}

func (r *Runner) Run(_ context.Context, t target.Template) (err error) {
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
	return
}

// check if already exist
func (r *Runner) checkTargetDir(dir, name string) (err error) {
	if _, err := os.Stat(dir); err == nil {
		err = fmt.Errorf("cannot create project %s: %s already exists", name, dir)
	} else if !os.IsNotExist(err) {
		err = fmt.Errorf("cannot create project %s: %v", name, err)
	}
	return
}

func (r *Runner) targetArgs(t target.Template) (args template.Args) {
	name := r.templates[t.Name()]
	args = r.args
	args.ProjectName = name
	args.PackageName = "new.portal." + name
	return
}
