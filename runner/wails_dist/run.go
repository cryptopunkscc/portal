package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/watcher"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/target"
	"path/filepath"
)

type Runner struct {
	send  target.MsgSend
	inner target.Runner[target.AppHtml]
}

func NewRunner(newRuntime bind.NewRuntime, send target.MsgSend) target.Runner[target.DistHtml] {
	return &Runner{
		send:  send,
		inner: wails.NewRunner(newRuntime),
	}
}

func (r *Runner) Reload() (err error) {
	return r.inner.Reload()
}

func (r *Runner) Run(ctx context.Context, dist target.DistHtml) (err error) {
	if !filepath.IsAbs(dist.Abs()) {
		return plog.Errorf("Runner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r)

	go func() {
		pkg := dist.Manifest().Package
		watch := watcher.Runner[target.DistHtml](func() (err error) {
			if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
				log.F().Println(err)
			}
			err = r.inner.Reload()
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

	if err = r.inner.Run(ctx, dist); err != nil {
		return
	}

	return
}
