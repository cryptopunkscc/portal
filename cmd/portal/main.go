package main

import (
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/runner/cli"
)

func main() {
	tokens := token.Repository{Dir: env.PortaldTokens.MkdirAll()}
	if t, err := tokens.Get("portal"); err == nil {
		apphost.Default.AuthToken = string(t.Token)
	}

	a := Application{
		Portal: apphost.Default.Portald(),
	}

	cli.Run(a.Handler())
}

type Application struct {
	Portal apphost.Portald
}
