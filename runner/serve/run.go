package serve

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"log"
)

func Run(ctx context.Context, port string, handlers rpc.Handlers, tray runtime.Tray) (err error) {
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), "service ", 0))
	for name, h := range handlers {
		s.RouteFunc(name, h)
	}

	go func() {
		if err = s.Run(ctx); err != nil {
			return
		}
	}()

	hasTray := tray != nil
	log.Printf("portal service started tray:%v", hasTray)

	if hasTray {
		tray(ctx)
	}
	<-ctx.Done()
	return
}
