package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target/dev/reload"
	"github.com/cryptopunkscc/portal/target/dist"
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

func (r *ReRunner) Run(ctx context.Context, distHtml target.DistHtml, args ...string) (err error) {
	if !filepath.IsAbs(distHtml.Abs()) {
		return plog.Errorf("ReRunner needs absolute path: %s", distHtml.Abs())
	}
	log := plog.Get(ctx).Type(r)

	go func() {
		pkg := distHtml.Manifest().Package
		watch := dist.ReRunner[target.DistHtml](func(...string) (err error) {
			if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
				log.F().Println(err)
			}
			err = r.Reload()
			if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
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
	r.send = reload.Start(ctx, distHtml, r.Reload, r.Core)
	if err = r.ReRunner.Run(ctx, distHtml, args...); err != nil {
		return
	}

	return
}
