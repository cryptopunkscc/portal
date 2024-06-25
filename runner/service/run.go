package service

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
)

func NewRun(handlers rpc.Handlers) target.Dispatch {
	return NewRunner(handlers).Run
}

type Runner struct {
	handlers rpc.Handlers
}

func NewRunner(handlers rpc.Handlers) *Runner {
	return &Runner{handlers: handlers}
}

func (r Runner) Run(ctx context.Context, port string, _ ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("start port:%s", port)
	defer log.Printf("exit port:%s", port)
	app := rpc.NewApp(port)
	for name, h := range r.handlers {
		app.RouteFunc(name, h)
	}
	if err = app.Run(ctx); err != nil {
		return plog.Err(err)
	}
	return
}
