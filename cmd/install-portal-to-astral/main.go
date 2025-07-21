package main

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() {
	cli.Run(handler)
}

var handler = cmd.Handler{
	Name: "install-portal-to-astral",
	Desc: "Astrald & portal environment installer.",
	Func: run,
	Params: cmd.Params{
		{
			Type: "string",
			Desc: "Optional user name. When specified, the installed node will be assigned to a new user identity associated with the name. Otherwise, the installed node will be ready to claim by existing user.",
		},
	},
}

func run(username string) (err error) {
	if username != "" {
		return firstInstallation(username)
	}
	return nextInstallation()
}
