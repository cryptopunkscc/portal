package template

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/git"
	"github.com/cryptopunkscc/portal/resolve/template"
)

type Runner struct {
	dir  string
	args template.Args
}

func NewRunner(dir string) *Runner {
	r := &Runner{
		dir:  target.Abs(dir),
		args: defaultTemplateArgs(),
	}
	return r
}

func defaultTemplateArgs() template.Args {
	args := template.Args{}
	args.AuthorName = git.UserName(false)
	if args.AuthorName == "" {
		args.AuthorName = git.UserName(true)
	}
	args.AuthorEmail = git.UserEmail(false)
	if args.AuthorEmail == "" {
		args.AuthorEmail = git.UserEmail(true)
	}
	return args
}
