package main

import (
	"context"

	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/api/version"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/source"
)

func main() { cli.Run(handler()) }

func handler() cmd.Handler {
	return cmd.Handler{
		Func: runner().Run,
		Name: "html",
		Desc: "Start portal app in wails runner.",
		Params: cmd.Params{
			{Type: "string", Desc: "The app name, or app package name, or release bundle ID, or absolute path to app bundle, or absolute path to app directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Name},
		},
	}
}

func runner() Dispatcher {
	return Dispatcher{
		Runner: RunFirst,
		Provider: Provider[Runnable]{
			Priority: Priority{
				Match[Bundle_],
				Match[Dist_],
			},
			Repository: Repositories{
				source.Repository,
				bundle.Repository{Apphost: apphost.Default},
			},
			Resolve: Any[Runnable](
				wails.Runner(core).Try,
			),
		},
	}
}

func core(ctx context.Context, portal Portal_) (bind.Core, context.Context) {
	r, ctx := bind.NewFrontendCore(ctx, portal)
	return &Adapter{r}, ctx
}

type Adapter struct{ bind.Core }
