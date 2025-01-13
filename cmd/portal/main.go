package main

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	apphost2 "github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/start"
	"github.com/cryptopunkscc/portal/runner/version"
	runtime "github.com/cryptopunkscc/portal/runtime/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() { cli.Run(Application{}.Handler()) }

type Application struct{}

func (a Application) Handler() cmd.Handler {
	return cmd.Handler{
		Name: "portal",
		Desc: "Portal command line.",
		Func: start.Create(a).Run,
		Params: cmd.Params{
			{Name: "open o", Type: "bool", Desc: "Open portal tha app as background process without redirecting IO."},
			{Name: "query q", Type: "string", Desc: "Optional query to execute on invoked service"},
			{Name: "dev d", Type: "bool", Desc: "Development mode."},
			{Type: "string", Desc: "Application source. The source can be a app name, package name, app bundle path or app dir."},
			{Type: "...string", Desc: "Optional application arguments."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}
}
func (a Application) Apphost() apphost.Client { return apphost2.Basic }
func (a Application) Portal() portal.Client   { return runtime.Client("portal") }
