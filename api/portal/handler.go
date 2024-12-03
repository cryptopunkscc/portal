package portal

import "github.com/cryptopunkscc/portal/runtime/rpc2/cmd"

func Handler(service Service) cmd.Handler {
	return cmd.Handler{
		Name:   "portald",
		Desc:   "Portal daemon.",
		Params: nil,
		Sub: cmd.Handlers{
			{
				Func: ping,
				Name: "ping",
			},
			{
				Func: service.Open(),
				Name: "open",
				Desc: "Open portal app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Absolute path to app bundle or directory."},
					{Type: "string", Desc: "Optional command to run on opened app."},
				},
			},
			{
				Func: service.Shutdown(),
				Name: "close",
				Desc: "Shutdown portal environment and close all running apps.",
			},
		},
	}
}

func ping() {}
