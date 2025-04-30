package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/watcher"
	"github.com/cryptopunkscc/portal/target/html"
	"path/filepath"
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
	send    target.MsgSend
	newCore bind.NewCore
}

func (r *ReRunner) Reload() (err error) {
	return r.ReRunner.Reload()
}

func (r *ReRunner) Run(ctx context.Context, dist target.DistHtml, args ...string) (err error) {
	if !filepath.IsAbs(dist.Abs()) {
		return plog.Errorf("ReRunner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r)

	go func() {
		pkg := dist.Manifest().Package
		watch := watcher.ReRunner[target.DistHtml](func(...string) (err error) {
			if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
				log.F().Println(err)
			}
			err = r.Reload()
			if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
				log.F().Println(err)
			}
			return err
		})
		err := watch.Run(ctx, dist)
		if err != nil {
			log.F().Println(err)
		}
	}()

	r.Core, ctx = r.newCore(ctx, dist)
	r.send = reload.Start(ctx, dist, r.Reload, r.Core)
	if err = r.ReRunner.Run(ctx, dist, args...); err != nil {
		return
	}

	return
}
