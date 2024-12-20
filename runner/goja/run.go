package goja

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

type Runner struct {
	newRuntime bind.NewRuntime
	backend    *Backend
	app        target.AppJs
}

func NewRunner(newRuntime bind.NewRuntime) target.Runner[target.AppJs] {
	return &Runner{newRuntime: newRuntime}
}

func NewRun(newRuntime bind.NewRuntime) target.Run[target.AppJs] {
	return NewRunner(newRuntime).Run
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.app.Files())
}

func (r *Runner) Run(ctx context.Context, app target.AppJs) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", app, app.Abs())
	r.app = app
	r.backend = NewBackend(r.newRuntime(ctx, app))
	if err = r.Reload(); err != nil {
		return
	}
	<-ctx.Done()
	return
}
