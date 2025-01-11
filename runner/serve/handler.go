package serve

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

type Service interface {
	Open() target.Run[string]
	Shutdown() context.CancelFunc
	Observe() func(ctx context.Context, conn rpc.Conn) (err error)
}

func Handler(service Service) cmd.Handler {
	open := service.Open()
	return cmd.Handler{
		Name: "portald",
		Desc: "Portal daemon.",
		Func: "RouteAll",
		Sub: cmd.Handlers{
			{
				Func: func() any { return 0 },
				Name: "ping",
			},
			{
				Func: open,
				Name: "open o",
				Desc: "Open portal app.",
				Params: cmd.Params{
					{Type: "string", Desc: "Absolute path to app bundle or directory."},
					{Type: "...string", Desc: "Optional arguments."},
				},
			},
			{
				Func: func(ctx context.Context, conn rpc.Conn, s string, args ...string) (err error) {
					ctx = exec.WithReadWriter(ctx, conn)
					_ = open(ctx, s, args...)
					return rpc.Close
				},
				Name: "connect c",
				Desc: "Open portal app and redirect standard IO to client.",
				Params: cmd.Params{
					{Type: "string", Desc: "Absolute path to app bundle or directory."},
					{Type: "...string", Desc: "Optional arguments."},
				},
			},
			{
				Func: service.Shutdown(),
				Name: "close",
				Desc: "Shutdown portal environment and close all running apps.",
			},
			{
				Func: service.Observe(),
				Name: "observe",
			},
		},
	}
}
