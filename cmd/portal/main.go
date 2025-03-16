package main

import (
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/runner/cli"
)

func main() {
	cli.Run(Application{
		Portal: portald.NewClient(),
	}.Handler())
}

type Application struct {
	Portal portald.Client
}
