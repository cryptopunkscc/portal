package portald

import (
	"github.com/cryptopunkscc/portal/client/apphost"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runner/uninstall"
	"github.com/cryptopunkscc/portal/runtime/apps"
	"github.com/cryptopunkscc/portal/runtime/rpc/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func (s *Runner[T]) Handlers() cmd.Handlers {
	return cmd.Handlers{
		cli.Handler,
		cli.EncodingHandler,
		cli.StdHandler,
		{
			Func: s.Api,
			Name: "api",
			Desc: "Print API schema.",
		},
		{
			Func: s.Ping,
			Name: "ping",
		},
		{
			Func: Join,
			Name: "join",
		},
		{
			Func: s.Open().Start,
			Name: "open o",
			Desc: "Open portal app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Absolute path to app bundle or directory."},
				{Type: "...string", Desc: "Optional arguments."},
			},
		},
	}.Plus(s.publicHandlers()...)
}

func (s *Runner[T]) publicHandlers() cmd.Handlers {
	return cmd.Handlers{
		{
			Func: s.Connect,
			Name: "connect c",
			Desc: "Open portal app and redirect standard IO to client.",
			Params: cmd.Params{
				{Type: "string", Desc: "Absolute path to app bundle or directory."},
				{Type: "...string", Desc: "Optional arguments."},
			},
		},
		{
			Func: install.Token,
			Name: "token",
			Desc: "Create a new token or return existing one.",
			Params: cmd.Params{
				{Type: "string", Desc: "Token name."},
			},
		},
		{
			Func: apphost.NewClient().ListTokens,
			Name: "tokens",
			Desc: "List apphost tokens.",
			Params: cmd.Params{
				{Name: "format f", Type: "string", Desc: "Format [json, bin]"},
			},
		},
		{
			Func: install.Runner(apps.Dir).Run,
			Name: "install i",
			Desc: "Install app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Path to containing directory"},
			},
		},
		{
			Func: uninstall.Runner(apps.Source),
			Name: "uninstall d",
			Desc: "Uninstall app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Application name or package name"},
			},
		},
		{
			Func: s.ListApps,
			Name: "list l",
			Desc: "List installed apps.",
			Params: cmd.Params{
				{Name: "hidden h", Type: "boolean", Desc: "Include hidden apps."},
			},
			Sub: cmd.Handlers{{
				Func: s.ObserveApps,
				Name: "observe o",
				Desc: "Observe installed list apps.",
				Params: cmd.Params{
					{Name: "hidden h", Type: "boolean", Desc: "Include hidden apps."},
				},
			}},
		},
		{
			Func: s.Shutdown,
			Name: "close",
			Desc: "Shutdown portal environment and close all running apps.",
		},
	}
}
