package wails_dist

import (
	"context"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/dev"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/deprecated/wails"
	"github.com/cryptopunkscc/portal/target/dev/reload"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/html"
)

func Runner(newCore bind.NewCore) *target.SourceRunner[target.DistHtml] {
	return &target.SourceRunner[target.DistHtml]{
		Runner: &ReRunner{
			newCore: newCore,
		},
		Resolve: target.Any[target.DistHtml](
			target.Skip("node_modules"),
			target.Try(html.ResolveDist),
			target.Try(html.ResolveBundle),
		),
	}
}

type ReRunner struct {
	*wails.ReRunner
	send    dev.SendMsg
	newCore bind.NewCore
}

func (r *ReRunner) Reload(ctx context.Context) (err error) {
	return r.ReRunner.Reload(ctx)
}

func (r *ReRunner) Run(ctx context.Context, distHtml target.DistHtml, args ...string) (err error) {
	if !filepath.IsAbs(distHtml.Abs()) {
		return plog.Errorf("ReRunner needs absolute path: %s", distHtml.Abs())
	}
	log := plog.Get(ctx).Type(r)

	go func() {
		pkg := distHtml.Manifest().Package
		watch := dist.ReRunner[target.DistHtml](func(...string) (err error) {
			if err := r.send(dev.NewMsg(pkg, dev.Changed)); err != nil {
				log.F().Println(err)
			}
			err = r.Reload(ctx)
			if err := r.send(dev.NewMsg(pkg, dev.Refreshed)); err != nil {
				log.F().Println(err)
			}
			return err
		})
		err := watch.Run(ctx, distHtml)
		if err != nil {
			log.F().Println(err)
		}
	}()

	r.Core, ctx = r.newCore(ctx, distHtml)
	r.send = reload.Start(ctx, distHtml.Manifest().Package, r.Reload, r.Core)
	if err = r.ReRunner.Run(ctx, distHtml, args...); err != nil {
		return
	}

	return
}
