package goja

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Runner struct {
	newApi  target.NewApi
	prefix  []string
	backend *Backend
	app     target.AppBackend
}

func NewRunner(newApi target.NewApi, prefix ...string) *Runner {
	return &Runner{newApi: newApi, prefix: prefix}
}

func (r *Runner) Reload() (err error) {
	return r.backend.RunFs(r.app.Files())
}

func (r *Runner) Run(ctx context.Context, app target.AppBackend) (err error) {
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
