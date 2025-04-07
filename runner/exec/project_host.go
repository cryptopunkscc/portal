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
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/project"
	"github.com/cryptopunkscc/portal/resolve/source"
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
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.ProjectHost")

		host := schemaPrefix
		if host == nil {
			opt := apphost.PortaldOpenOpt{}
			if opt.Load(ctx); len(opt.Schema) > 0 {
				host = append(host, opt.Schema)
			}
		}

		manifest := src.Manifest()
		if manifest.Schema != "" {
			host = append(host, manifest.Schema)
		}

		hostId := strings.Join(host, ".")
		log.Println("running:", hostId, manifest.Package, args)

		runners, err := target.
			FindByPath(source.File, exec2.ResolveProject).
			OrById(path.Resolver(exec2.ResolveProject, env.PortaldApps.Source())).
			Call(ctx, hostId)

		if err != nil {
			return
		}
		if len(runners) == 0 {
			return target.ErrNotFound
		}
		var runner target.ProjectExec
		for _, r := range runners {
			if r.Manifest().Exec != "" {
				runner = r
				break
			}
		}
		if runner == nil {
			return target.ErrNotFound
		}

		e := runner.Manifest().Exec
		args = slices.Insert(args, 0, src.Abs())
		c, err := exec.Cmd{}.Parse(e, runner.Abs(), strings.Join(args, " "))
		if err != nil {
			return
		}
		log.Println("running", c)
		return Cmd{}.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
	}
}
