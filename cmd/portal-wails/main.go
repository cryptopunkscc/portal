package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/v2/wails"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Func: run,
	Name: "portal-wails",
	Desc: "Start portal HTML app in wails runner.",
	Params: cmd.Params{
		{Type: "string", Desc: "One of: app name, app package name, release bundle ID, absolute path to app bundle, absolute path to app directory."},
	},
	Sub: cmd.Handlers{
		{Name: "v", Desc: "Print version.", Func: version.Name},
	},
}

func run(ctx context.Context, src string) (err error) {
	s := source.Providers{
		source.OsFs,
		app.Objects{}.Default(),
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	adapter := &Adapter{}
	adapter.Core, ctx = bind.DefaultCoreFactory{}.Create(ctx)
	for _, ss := range source.Collect(s,
		wails.NewAppRunner(adapter),
		wails.NewBundleRunner(adapter),
	) {
		return ss.(wails.Runner).Run(ctx)
	}

	return fs.ErrInvalid
}

type Adapter struct{ bind.Core }
