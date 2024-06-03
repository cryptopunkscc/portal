package goja

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Runner struct {
	newApi  target.NewApi
	backend *Backend
	app     target.AppJs
}

func NewRunner(newApi target.NewApi) *Runner {
	return &Runner{newApi: newApi}
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.app.Files())
}

func (r *Runner) Run(ctx context.Context, app target.AppJs) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", app, app.Abs())
	r.app = app
	r.backend = NewBackend(r.newApi(ctx, app))
	if err = r.Reload(); err != nil {
		return
	}
	<-ctx.Done()
	return
}
