package portald

import (
	"github.com/cryptopunkscc/portal/client/apphost"
	"github.com/cryptopunkscc/portal/runner/install"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
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
			Func: s.Open().Start,
			Name: "open o",
			Desc: "Open portal app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Absolute path to app bundle or directory."},
				{Type: "...string", Desc: "Optional arguments."},
			},
		},
		{
			Func: s.Connect,
			Name: "connect c",
			Desc: "Open portal app and redirect standard IO to client.",
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
			Func: apphost.NewClient().ListTokens,
			Name: "tokens",
			Desc: "List apphost tokens.",
			Params: cmd.Params{
				{Name: "format f", Type: "string", Desc: "Format [json, bin]"},
			},
		},
		{
			Func: install.Token,
			Name: "token",
			Desc: "Create a new token.",
			Params: cmd.Params{
				{Type: "string", Desc: "Token name."},
			},
		},
		{
			Func: s.Install,
			Name: "install i",
			Desc: "Install app.",
			Params: cmd.Params{
				{Type: "string", Desc: "Path to containing directory"},
			},
		},
		{
			Func: s.Uninstall,
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
		},
		{
			Func: s.Shutdown,
			Name: "close",
			Desc: "Shutdown portal environment and close all running apps.",
		},
	}
}
