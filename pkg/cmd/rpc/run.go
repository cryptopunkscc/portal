package rpc

import (
	"context"
	jrpc "github.com/cryptopunkscc/go-apphost-jrpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
)

func Run(bindings runner.Bindings) (err error) {

	s := jrpc.NewApp("portal")
	s.Logger(log.New(log.Writer(), "service ", 0))
	s.With(bindings)
	s.RouteFunc("open", portal.Open)
	s.RouteFunc("observe", appstore.Observe)
	s.RouteFunc("install", appstore.Install)
	s.RouteFunc("uninstall", appstore.Uninstall)

	ctx := context.Background()
	if err = s.Run(ctx); err != nil {
		return
	}
	<-ctx.Done()
	return
}
