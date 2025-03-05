package main

import (
	"github.com/cryptopunkscc/portal/runner/version"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func (a Application) Handler() (h cmd.Handler) {
	h = cmd.Handler{
		Name: "Portal CLI",
		Desc: "Manage and run portal apps.",
		Func: a.Run,
		Params: cmd.Params{
			{Name: "open o", Type: "bool", Desc: "Open portal tha app as background process without redirecting IO."},
			{Name: "query q", Type: "string", Desc: "Optional query to execute on invoked service"},
			{Name: "dev d", Type: "bool", Desc: "Development mode."},
			{Type: "string", Desc: "Application source. The source can be a app name, package name, app bundle path or app dir."},
			{Type: "...string", Desc: "Optional application arguments."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
			{
				Func: a.Arg,
				Name: "-a",
				Desc: "Execute list of commands with given arg.",
				Params: cmd.Params{
					{Type: "string", Desc: "Argument value."},
					{Type: "...string", Desc: "List of commands to run with given arg."},
				},
			},
		},
	}
	cmd.InjectHelp(&h)
	a.injectPortaldApi(&h)
	return
}
