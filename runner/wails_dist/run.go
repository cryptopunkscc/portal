package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/watcher"
	"github.com/cryptopunkscc/go-astral-js/target"
	"path"
)

type Runner struct {
	send  target.MsgSend
	inner *wails.Runner
}

func NewRunner(newApi target.NewApi, send target.MsgSend) (runner *Runner) {
	runner = &Runner{
		send:  send,
		inner: wails.NewRunner(newApi),
	}
	return
}

func (r *Runner) Reload() (err error) {
	return r.inner.Reload()
}

func (r *Runner) Run(ctx context.Context, dist target.DistHtml) (err error) {
	if !path.IsAbs(dist.Abs()) {
		return plog.Errorf("Runner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r)

	go func() {
		pkg := dist.Manifest().Package
		watch := watcher.NewRunner[target.DistHtml](func() (err error) {
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
