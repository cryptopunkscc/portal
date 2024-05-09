package serve

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"log"
)

func Run(ctx context.Context, bindings runtime.New, handlers rpc.Handlers, tray runtime.Tray) (err error) {
	s := rpc.NewApp("portal")
	s.Logger(log.New(log.Writer(), "service ", 0))
	s.With(bindings)
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
