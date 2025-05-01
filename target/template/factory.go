package template

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/git"
)

type Factory struct {
	dir  string
	args Args
}

// Args will be embedded into the tmpl files during the installation
type Args struct {
	ProjectName string
	PackageName string
	AuthorName  string
	AuthorEmail string
	Description string
	Url         string
}

func ProjectFactory(dir string) *Factory {
	r := &Factory{
		dir:  target.Abs(dir),
		args: defaultTemplateArgs(),
	}
	return r
}

func defaultTemplateArgs() Args {
	args := Args{}
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
