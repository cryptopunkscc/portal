package jrpc

import (
	"context"
	jrpc "github.com/cryptopunkscc/go-apphost-jrpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/build"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/create"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/dev"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/prod"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/publish"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"log"
)

func Run(bindings runner.Bindings) (err error) {

	s := jrpc.NewApp("portal")
	s.Logger(log.New(log.Writer(), "service ", 0))

	s.With(bindings)
	s.Routes("*")
	s.RouteFunc("open", prod.Run)
	s.RouteFunc("create", create.Run)
	s.RouteFunc("dev", dev.Run)
	s.RouteFunc("build", build.Run)
	s.RouteFunc("bundle", bundle.Run)
	s.RouteFunc("publish", publish.Run)

	ctx := context.Background()
	err = s.Run(ctx)
	if err != nil {
		return
	}
	<-ctx.Done()
	return
}
