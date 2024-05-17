package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"log"
)

func Run(ctx context.Context, port string, handlers rpc.Handlers) (err error) {
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), "service ", 0))
	for name, h := range handlers {
		s.RouteFunc(name, h)
	}

	if err = s.Run(ctx); err != nil {
		return fmt.Errorf("serve.Run exit: %w", err)
	}

	return
}
