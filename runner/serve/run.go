package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"log"
)

type Runner struct {
	handlers rpc.Handlers
}

func NewRunner(handlers rpc.Handlers) *Runner {
	return &Runner{handlers: handlers}
}

func (r Runner) Run(ctx context.Context, port string, _ ...string) (err error) {
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), "service ", 0))
	for name, h := range r.handlers {
		s.RouteFunc(name, h)
	}

	if err = s.Run(ctx); err != nil {
		return fmt.Errorf("serve.Run exit: %w", err)
	}

	return
}
