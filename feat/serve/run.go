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
	s.RouteFunc("ping", func() {})
	s.RouteFunc("open", portal.CmdOpenerCtx(ctx))
	s.RouteFunc("observe", appstore.Observe)
	s.RouteFunc("install", appstore.Install)
	s.RouteFunc("uninstall", appstore.Uninstall)

	go func() {
		if err = s.Run(ctx); err != nil {
			return
		}
	}()

	log.Printf("portal service started tray:%v", indicator)

	if indicator {
		tray.Run(ctx)
	}
	<-ctx.Done()
	return
}
