package exec

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/dist"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/project"
	"slices"
	"strings"
)

var ProjectHostRunner = target.SourceRunner[target.Portal_]{
	Resolve: target.Any[target.Portal_](
		target.Try(dist.ResolveAny),
		target.Try(project.ResolveAny),
	),
	Runner: &projectHostRunner{},
}

type projectHostRunner struct{}

func (r *projectHostRunner) Run(ctx context.Context, src target.Portal_, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if src.Manifest().Schema == "" {
		return errors.New("projectHostRunner requires a schema declared in manifest")
	}

	log := plog.Get(ctx).Type(r)
	repo := target.SourcesRepository[target.ProjectExec]{
		Sources: []target.Source{env.PortaldApps.Source()},
		Resolve: target.Any[target.ProjectExec](exec2.ResolveProject.Try),
	}
	hostId := src.Manifest().Schema
	opt := apphost.PortaldOpenOpt{}
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

	return Cmd{}.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
}
