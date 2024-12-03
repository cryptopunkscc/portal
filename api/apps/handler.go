package apps

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

type Service interface {
	List() func() []target.App_
	Install() func(src string) (err error)
	Uninstall() func(id string) (err error)
}

func Handler(service Service) cmd.Handler {
	return cmd.Handler{
		Name: "portal-apps",
		Desc: "Portal applications management.",
		Sub: cmd.Handlers{
			{
				Name: "list l",
				Func: service.List(),
			},
			{
				Name: "install i",
				Desc: "Install app from a given portal app bundle path.",
				Func: service.Install(),
				Params: cmd.Params{
					{Type: "string", Desc: "Path to containing directory"},
				},
			},
			{
				Name: "delete d",
				Desc: "Uninstall app.",
				Func: service.Uninstall(),
				Params: cmd.Params{
					{Type: "string", Desc: "Application name or package name"},
				},
			},
		},
	}
}
