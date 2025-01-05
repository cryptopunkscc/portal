package main

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
	apphost2 "github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/start"
	runtime "github.com/cryptopunkscc/portal/runtime/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	cli.Run(cmd.Handler{
		Name: "portal",
		Desc: "Portal command line.",
		Func: start.Create(deps{}).Run,
		Params: cmd.Params{
			{Name: "query q", Type: "string", Desc: "Optional query to execute on invoked service"},
			{Type: "string", Desc: "Application source. The source can be a app name, package name, app bundle path or app dir."},
			{Type: "...string", Desc: "Optional application arguments."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	})
}

type deps struct{}

func (m deps) Apphost() apphost.Client { return apphost2.Basic }
func (m deps) Portal() portal.Client   { return runtime.Client("portal") }
