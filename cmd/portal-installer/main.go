package main

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() {
	cli.Run(handler)
}

var handler = cmd.Handler{
	Desc: "Astrald & portal environment installer.",
	Sub: cmd.Handlers{
		{
			Func: firstInstallation,
			Name: "first",
			Desc: "First installation. Choose if you are installing for the first time. It will install complete environment with default apps and generate your identity.",
			Params: cmd.Params{
				{
					Type: "string",
					Desc: "Your user name that will be associated with generated identity.",
				},
			},
		},
		{
			Func: nextInstallation,
			Name: "next",
			Desc: "Next installation. Choose if you already have your identity generated but want to claim another device. It will install base environment ready to claim from another device.",
		},
	},
}
