package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/watcher"
	"github.com/cryptopunkscc/go-astral-js/target"
	"path"
)

type Runner struct {
	newApi  target.NewApi
	send    target.MsgSend
	dist    target.DistBackend
	backend *goja.Backend
}

func NewRunner(newApi target.NewApi, send target.MsgSend) *Runner {
	return &Runner{newApi: newApi, send: send}
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.dist.Files())
}

func (r *Runner) Run(ctx context.Context, dist target.DistBackend) (err error) {
	if !path.IsAbs(dist.Abs()) {
		return plog.Errorf("Runner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", dist, dist.Abs())
	r.dist = dist
	r.backend = goja.NewBackend(r.newApi(ctx, dist))
	if err = r.Reload(); err != nil {
		return
	}
	pkg := dist.Manifest().Package
	watch := watcher.NewRunner[target.DistBackend](func() (err error) {
		if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
			log.F().Println(err)
		}
		err = r.Reload()
		if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
			log.F().Println(err)
		}
		return err
	})
	return watch.Run(ctx, dist)
}
