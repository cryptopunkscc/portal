package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/watcher"
	"path/filepath"
	"time"
)

type Runner struct {
	newCore bind.NewCore
	send    target.MsgSend
	dist    target.DistJs
	backend *goja.Backend
}

func NewRunner(newCore bind.NewCore, send target.MsgSend) target.ReRunner[target.DistJs] {
	return &Runner{newCore: newCore, send: send}
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.dist.FS())
}

func (r *Runner) Run(ctx context.Context, dist target.DistJs, args ...string) (err error) {
	if any(r.newCore) == nil {
		panic("newCore cannot be nil")
	}
	if !filepath.IsAbs(dist.Abs()) {
		return plog.Errorf("ReRunner needs absolute path: %s", dist.Abs())
	}
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", dist, dist.Abs())
	core, ctx := r.newCore(ctx, dist)
	r.backend = goja.NewBackend(core)
	r.dist = dist
	if err = r.Reload(); err != nil {
		log.E().Println(err.Error())
	}
	pkg := dist.Manifest().Package
	watch := watcher.ReRunner[target.DistJs](func(...string) error {
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
	return watch.Run(ctx, dist, args...)
}
