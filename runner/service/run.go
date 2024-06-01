package service

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
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

func (r Runner) Start(ctx context.Context, port string, args ...string) {
	go func() {
		if err := r.Run(ctx, port, args...); err != nil {
			plog.Get(ctx).Type(r).F().Printf("%s: %v", port, err)
		}
	}()
}

func (r Runner) Run(ctx context.Context, port string, _ ...string) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)
	app := rpc.NewApp(port)
	for name, h := range r.handlers {
		app.RouteFunc(name, h)
	}
	if err = app.Run(ctx); err != nil {
		return plog.Err(err)
	}
	return
}
