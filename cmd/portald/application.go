package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

type Application struct{ portald.Service }

func (a *Application) commands() (h cmd.Handler) {
	return cmd.Handler{
		Func: a.run,
		Name: "portald",
		Desc: "Portal daemon.",
		Params: []cmd.Param{
			{
				Name: "config c",
				Type: "string",
				Desc: "Path to the config file or directory containing .portal.yml",
			},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}
