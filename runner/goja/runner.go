package goja

import (
	"context"
	"fmt"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/js"
)

func Runner(newCore bind.NewCore) *target.SourceRunner[target.AppJs] {
	return &target.SourceRunner[target.AppJs]{
		Resolve: target.Any[target.AppJs](
			js.ResolveDist.Try,
			js.ResolveBundle.Try,
		),
		Runner: NewRunner(newCore),
	}
}

type runner struct {
	newCore bind.NewCore
	backend *Backend
	app     target.AppJs
	args    []string
}

func NewRunner(newCore bind.NewCore) target.ReRunner[target.AppJs] {
	return &runner{newCore: newCore}
}

func (r *runner) Reload() (err error) {
	return r.backend.RunFs(r.app.FS(), r.args...)
}

func (r *runner) Run(ctx context.Context, app target.AppJs, args ...string) (err error) {
	// TODO pass args to js
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", app, app.Abs())
	core, ctx := r.newCore(ctx, app)
	r.app = app
	r.args = args
	r.backend = NewBackend(core)
	if err = r.Reload(); err != nil {
		return
	}
	<-ctx.Done()
	r.backend.Interrupt()
	if core.Code() > 0 {
		err = fmt.Errorf("exit %d", core.Code())
	}
	return
}
