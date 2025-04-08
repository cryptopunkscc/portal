package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"slices"
	"strings"
)

var BundleHostRunner = target.SourceRunner[target.Portal_]{
	Resolve: target.Any[target.Portal_](target.Try(bundle.ResolveAny)),
	Runner:  BundleHost(),
}

func BundleHost(schemaPrefix ...string) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		repo := target.SourcesRepository[target.BundleExec]{
			Sources: []target.Source{env.PortaldApps.Source()},
			Resolve: target.Any[target.BundleExec](exec.ResolveBundle.Try),
		}
		defer plog.TraceErr(&err)
		log := plog.Get(ctx).Scope("exec.BundleHost")
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
		log.Println("running:", hostId, src.Manifest().Package, args)

		hostBundle := repo.First(hostId)
		if hostBundle == nil {
			return target.ErrNotFound
		}

		hostExec, err := unpackExecutable(hostBundle)
		if err != nil {
			return
		}

		args = slices.Insert(args, 0, src.Abs())
		return Cmd{}.RunApp(ctx, *src.Manifest(), hostExec.Name(), args...)
	}
}
