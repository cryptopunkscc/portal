package exec

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/bundle"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/source"
	"slices"
)

func (r Runner) BundleHost() *target.SourceRunner[target.Portal_] {
	return &target.SourceRunner[target.Portal_]{
		Resolve: target.Any[target.Portal_](target.Try(bundle.ResolveAny)),
		Runner:  &BundleHostRunner{BundleRunner{r}},
	}
}

type BundleHostRunner struct{ BundleRunner }

func (r *BundleHostRunner) Run(ctx context.Context, src target.Portal_, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if src.Manifest().Schema == "" {
		return errors.New("BundleHostRunner requires a schema declared in manifest")
	}

	log := plog.Get(ctx).Type(r)
	repo := target.SourcesRepository[target.BundleExec]{
		Sources: []target.Source{source.Dir(r.Apps)},
		Resolve: target.Any[target.BundleExec](exec.ResolveBundle.Try),
	}
	hostId := src.Manifest().Schema
	opt := apphost.PortaldOpenOpt{}
	if opt.Load(ctx); len(opt.Schema) > 0 {
		hostId = hostId + "." + opt.Schema
	}

	log.Println("running:", hostId, src.Manifest().Package, args)

	hostBundle := repo.First(hostId)
	if hostBundle == nil {
		return target.ErrNotFound
	}

	hostExec, err := r.unpackExecutable(hostBundle)
	if err != nil {
		return
	}

	args = slices.Insert(args, 0, src.Abs())
	return r.RunApp(ctx, *src.Manifest(), hostExec.Name(), args...)
}
