package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

type Service interface {
	Open() target.Request
	Shutdown() context.CancelFunc
}

func Handler(service Service) cmd.Handler {
	return cmd.Handler{
		Name: "portald",
		Desc: "Portal daemon.",
		Sub: cmd.Handlers{
			{
				Func: func() {},
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
