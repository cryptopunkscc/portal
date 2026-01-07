package main

import (
	"context"
	"io/fs"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/cryptopunkscc/portal/runner/v2/goja"
	goja_dist "github.com/cryptopunkscc/portal/runner/v2/goja/dist"
	goja_pro "github.com/cryptopunkscc/portal/runner/v2/goja/pro"
	"github.com/cryptopunkscc/portal/runner/v2/wails"
	wails_dist "github.com/cryptopunkscc/portal/runner/v2/wails/dist"
	wails_pro "github.com/cryptopunkscc/portal/runner/v2/wails/pro"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/npm"
	"github.com/cryptopunkscc/portal/source/tmpl"
)

func main() { cli.Run(handler) }

var handler = cmd.Handler{
	Name: "portal-dev",
	Desc: "Development kit for Astral apps",
	Sub: cmd.Handlers{
		cmd.Handler{
			Func: tmpl.Create,
			Name: "create c",
			Desc: "Create new Astral app from template",
		},
		cmd.Handler{
			Func: npm.BuildNpmApps,
			Name: "build b",
			Desc: "Build Astral apps",
		},
		cmd.Handler{
			Func: app.Publisher{}.PublishBundles,
			Name: "publish p",
			Desc: "Publish app bundles to Astral",
		},
		cmd.Handler{
			Func: runApp,
			Name: "run r",
			Desc: "Run HTML/JS app in hot reloading runner.",
		},
	},
}

func runApp(ctx context.Context, src string, args ...string) (err error) {
	s := source.Providers{
		source.OsFs,
	}.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}

	core, ctx := bind.DefaultCoreFactory{}.Create(ctx)
	for _, ss := range source.Collect(s,
		goja_pro.NewRunner(core),
		goja_dist.NewRunner(core),
		goja.NewBundleRunner(core),
		wails_pro.NewRunner(core),
		wails_dist.NewRunner(core),
		wails.NewBundleRunner(core),
	) {
		switch r := ss.(type) {
		case goja.Runner:
			return r.Run(ctx, args...)
		case wails.Runner:
			return r.Run(ctx)
		}
	}

	return fs.ErrInvalid
}
