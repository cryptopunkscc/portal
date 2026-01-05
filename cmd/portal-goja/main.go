package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/v2/goja"
	"github.com/cryptopunkscc/portal/source"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Func: run,
	Name: "portal-goja",
	Desc: "Start portal JS app in goja runner.",
	Params: cmd.Params{
		{Type: "string", Desc: "One of: app name, app package name, release bundle ID, absolute path to app bundle, absolute path to app directory."},
	},
	Sub: cmd.Handlers{
		{Name: "v", Desc: "Print version.", Func: version.Name},
	},
}

func run(ctx context.Context, src string, args ...string) (err error) {
	s := source.Providers{
		source.OsFs,
		//app.Objects{Client: *astrald.DefaultClient()},
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	core, ctx := bind.DefaultCoreFactory{}.Create(ctx)
	for _, ss := range source.Collect(s,
		goja.NewAppRunner(core),
		goja.NewBundleRunner(core),
	) {
		return ss.(goja.Runner).Run(ctx, args...)
	}

	return fs.ErrInvalid
}
