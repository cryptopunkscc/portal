package exec

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/dist"
	exec2 "github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/source"
	"slices"
	"strings"
)

func (r Runner) ProjectHost() *target.SourceRunner[target.Portal_] {
	return &target.SourceRunner[target.Portal_]{
		Resolve: target.Any[target.Portal_](
			target.Try(dist.Resolve_),
			target.Try(project.Resolve_),
		),
		Runner: &ProjectHostRunner{r},
	}
}

type ProjectHostRunner struct{ Runner }

func (r *ProjectHostRunner) Run(ctx context.Context, src target.Portal_, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if src.Manifest().Schema == "" {
		return errors.New("ProjectHostRunner requires a schema declared in manifest")
	}

	log := plog.Get(ctx).Type(r)
	repo := target.SourcesRepository[target.Portal_]{
		Sources: []target.Source{source.Dir(r.Apps)},
		Resolve: target.Any[target.Portal_](exec2.ResolveProject.Try),
	}
	hostId := src.Manifest().Schema
	opt := portald.OpenOpt{}
	if opt.Load(ctx); len(opt.Schema) > 0 {
		hostId = hostId + "." + opt.Schema
	}

	log.Println("running:", hostId, src.Manifest(), args)

	hostBundle := repo.First(hostId)
	if hostBundle == nil || hostBundle.Manifest().Exec == "" {
		return target.ErrNotFound
	}

	hostExec := hostBundle.Manifest().Exec
	args = slices.Insert(args, 0, src.Abs())

	c, err := exec.Cmd{}.Parse(hostExec, hostBundle.Abs(), strings.Join(args, " "))
	if err != nil {
		return
	}

	return r.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
}
