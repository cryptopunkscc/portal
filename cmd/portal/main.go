package main

import (
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/runner/cli"
)

func main() {
	cli.Run(Application{
		Portal: apphost.Default.Portald(),
	}.Handler())
}

type Application struct {
	Portal apphost.Portald
}
