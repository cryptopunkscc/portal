package serve

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/tray"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"log"
)

func Run(ctx context.Context, bindings runtime.New, indicator bool) (err error) {
	s := rpc.NewApp("portal")
	s.Logger(log.New(log.Writer(), "service ", 0))
	s.With(bindings)
	s.RouteFunc("open", portal.OpenWithContext(ctx))
	s.RouteFunc("observe", appstore.Observe)
	s.RouteFunc("install", appstore.Install)
	s.RouteFunc("uninstall", appstore.Uninstall)

	if err = s.Run(ctx); err != nil {
		return
	}
	log.Println("portal service started")

	if indicator {
		tray.Run(ctx)
	}
	<-ctx.Done()
	return
}
