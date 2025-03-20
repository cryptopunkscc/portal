package goja

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Runner struct {
	newCore bind.NewCore
	backend *Backend
	app     target.AppJs
	args    []string
}

func NewRunner(newCore bind.NewCore) target.ReRunner[target.AppJs] {
	return &Runner{newCore: newCore}
}

func NewRun(newCore bind.NewCore) target.Run[target.AppJs] {
	return NewRunner(newCore).Run
}

func (r *Runner) ReRun() (err error) {
	return r.backend.RunFs(r.app.FS(), r.args...)
}

func (r *Runner) Run(ctx context.Context, app target.AppJs, args ...string) (err error) {
	// TODO pass args to js
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", app, app.Abs())
	core, ctx := r.newCore(ctx, app)
	r.app = app
	r.args = args
	r.backend = NewBackend(core)
	if err = r.ReRun(); err != nil {
		return
	}
	<-ctx.Done()
	return
}
