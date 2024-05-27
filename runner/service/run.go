package service

import (
	"context"
	"fmt"
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

func (r Runner) Run(ctx context.Context, port string, _ ...string) (err error) {
	plog.Get(ctx).Type(r).Set(&ctx)
	s := rpc.NewApp(port)
	//s.Logger(log.New(log.Writer(), "service ", 0))
	for name, h := range r.handlers {
		s.RouteFunc(name, h)
	}

	if err = s.Run(ctx); err != nil {
		return fmt.Errorf("serve.Run exit: %w", err)
	}

	return
}
