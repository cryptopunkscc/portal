package main

import (
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

func main() {
	cli.Run(Application{
		Portal: apphost.Default.Portald(),
	}.Handler())
}

type Application struct {
	Portal apphost.Portald
}
