package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/runner/cli"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
)

func main() {
	cli.Run(Application{
		Connect: apphost.Connect,
		Apphost: apphost.DefaultCached,
		Portal:  portald.NewClient(),
	}.Handler())
}

type Application struct {
	Connect func(context.Context) error
	Portal  portald.Client
	Apphost apphost.Client
}
