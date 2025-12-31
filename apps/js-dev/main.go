package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/deprecated/goja"
	"github.com/cryptopunkscc/portal/runner/deprecated/goja/dist"
	"github.com/cryptopunkscc/portal/runner/deprecated/goja/pro"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(Application{}.handler()) }

type Application struct{}

func (a Application) handler() cmd.Handler {
	return cmd.Handler{
		Func: source.File.NewRun(
			goja_pro.Runner(bind.CreateCore).Try,
			goja_dist.Runner(bind.CreateCore).Try,
			goja.Runner(bind.CreateCore).Try,
		),
		Name: "dev-js",
		Desc: "Start portal js app development in goja runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
			{
				Func: js.RunCreate,
				Name: "new n",
				Desc: "Create a new js app.",
				Params: cmd.Params{
					{Name: "template t", Type: "string", Desc: "Template to use."},
					{Type: "string", Desc: "Project destination directory."},
				},
			},
			{
				Func: js.ListTemplates,
				Name: "templates t",
				Desc: "List available templates.",
			},
		},
	}
}
