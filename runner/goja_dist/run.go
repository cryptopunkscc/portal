package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/broadcast"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/watcher"
	"github.com/cryptopunkscc/go-astral-js/target"
	"path"
)

type Runner struct {
	ctrlPort string
	newApi   target.NewApi
	backend  *goja.Backend
	dist     target.DistBackend
}

func NewRunner(ctrlPort string, newApi target.NewApi) *Runner {
	return &Runner{ctrlPort: ctrlPort, newApi: newApi}
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
	if err = r.backend.RunFs(dist.Files()); err != nil {
		return
	}
	pkg := dist.Manifest().Package
	watch := watcher.NewRunner[target.DistBackend](func() (err error) {
		err = broadcast.Send(r.ctrlPort, broadcast.NewMsg(pkg, broadcast.Changed))
		if err != nil {
			log.Println("broadcast.Send:", err)
		}
		err = r.backend.RunFs(dist.Files())
		err = broadcast.Send(r.ctrlPort, broadcast.NewMsg(pkg, broadcast.Refreshed))
		if err != nil {
			log.Println("broadcast.Send:", err)
		}
		return err
	})
	return watch.Run(ctx, dist)
}
