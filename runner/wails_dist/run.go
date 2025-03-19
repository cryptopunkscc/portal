package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/watcher"
	"path/filepath"
)

type reRunner struct {
	send  target.MsgSend
	inner target.ReRunner[target.AppHtml]
}

func ReRunner(newRuntime bind.NewRuntime, send target.MsgSend) target.ReRunner[target.DistHtml] {
	return &reRunner{
		send:  send,
		inner: wails.ReRunner(newRuntime),
	}
}

func (r *reRunner) ReRun() (err error) {
	return r.inner.ReRun()
}

func (r *reRunner) Run(ctx context.Context, dist target.DistHtml, args ...string) (err error) {
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
			err = r.inner.ReRun()
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

	if err = r.inner.Run(ctx, dist, args...); err != nil {
		return
	}

	return
}
