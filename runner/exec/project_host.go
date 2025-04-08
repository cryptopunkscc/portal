package exec

import (
	"context"
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
	Runner: ProjectHost(),
}

func ProjectHost(schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		repo := target.SourcesRepository[target.ProjectExec]{
			Sources: []target.Source{env.PortaldApps.Source()},
			Resolve: target.Any[target.ProjectExec](exec2.ResolveProject.Try),
		}
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.ProjectHost")
		host := schemaPrefix
		if host == nil {
			opt := apphost.PortaldOpenOpt{}
			if opt.Load(ctx); len(opt.Schema) > 0 {
				host = append(host, opt.Schema)
			}
		}
		if src.Manifest().Schema != "" {
			host = append(host, src.Manifest().Schema)
		}
		hostId := strings.Join(host, ".")
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

		log.Println("running", c)
		return Cmd{}.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
	}
}
