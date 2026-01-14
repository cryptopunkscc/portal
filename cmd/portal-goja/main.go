package main

import (
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/cmd/portal-goja/src"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Func: portal_goja.Application{Adapter: apphost.Default}.Run,
	Name: "portal-goja",
	Desc: "Start portal JS app in goja runner.",
	Params: cmd.Params{
		{Type: "string", Desc: "One of: app name, app package name, release bundle ID, absolute path to app bundle, absolute path to app directory."},
	},
	Sub: cmd.Handlers{
		{Name: "v", Desc: "Print version.", Func: version.Name},
	},
}
