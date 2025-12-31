package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/os"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/wails"
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
	src = os.Abs(src)
	s := source.Providers{
		source.OsFs,
		app.Objects{}.Default(),
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	adapter := &Adapter{}
	adapter.Core = bind.DefaultCoreFactory{}.Create(ctx)
	for _, ss := range source.Collect(s,
		&wails.AppRunner{},
		&wails.BundleRunner{},
	) {
		return ss.(wails.Runner).Run(adapter)
	}

	return fs.ErrInvalid
}

type Adapter struct{ bind.Core }
