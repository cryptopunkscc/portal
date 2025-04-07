package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"slices"
	"strings"
)

var BundleHostRunner = target.SourceRunner[target.Portal_]{
	Resolve: target.Any[target.Portal_](target.Try(bundle.ResolveAny)),
	Runner:  BundleHost(),
}

func BundleHost(schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.BundleHost")

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
			FindByPath(source.File, exec.ResolveBundle).
			OrById(path.Resolver(exec.ResolveBundle, env.PortaldApps.Source())).
			Call(ctx, hostId)

		if err != nil {
			return
		}
		if len(runners) == 0 {
			return target.ErrNotFound
		}

		var runner = runners[0]
		runnerExecutable, err := unpackExecutable(runner)
		if err != nil {
			return
		}

		args = slices.Insert(args, 0, src.Abs())
		return Cmd{}.RunApp(ctx, *src.Manifest(), runnerExecutable.Name(), args...)
	}
}
