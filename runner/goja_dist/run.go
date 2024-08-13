package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/watcher"
	"github.com/cryptopunkscc/portal/target"
	"path/filepath"
	"time"
)

type Runner struct {
	newRuntime target.NewRuntime
	send       target.MsgSend
	dist       target.DistJs
	backend    *goja.Backend
}

func NewRunner(newRuntime target.NewRuntime, send target.MsgSend) target.Runner[target.DistJs] {
	return &Runner{newRuntime: newRuntime, send: send}
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.dist.Files())
}

func (r *Runner) Run(ctx context.Context, dist target.DistJs) (err error) {
	if !filepath.IsAbs(dist.Abs()) {
		return plog.Errorf("Runner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", dist, dist.Abs())
	r.dist = dist
	if any(r.newRuntime) == nil {
		panic("newRuntime cannot be nil")
	}
	r.backend = goja.NewBackend(r.newRuntime(ctx, dist))
	if err = r.Reload(); err != nil {
		log.E().Println(err.Error())
	}
	pkg := dist.Manifest().Package
	watch := watcher.NewRunner[target.DistJs](func() error {
		if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
			log.E().Println(err)
		}
		if err := r.Reload(); err != nil {
			log.E().Println(err.Error())
		}
		// TODO find better solution then sleep
		// target.DevRefreshed msg must be delayed until backend is fully refreshed (all ports registered).
		time.Sleep(2 * time.Second)
		if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
			log.E().Println(err)
		}
		return nil
	})
	return watch.Run(ctx, dist)
}
