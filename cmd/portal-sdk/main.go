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
	"github.com/cryptopunkscc/portal/source/html"
	"github.com/cryptopunkscc/portal/source/js"
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
		cmd.Handler{
			Func: listTargets,
			Name: "list l",
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

	for _, ss := range source.Collect(s,
		&goja_pro.Runner{},
		&goja_dist.Runner{},
		&goja.BundleRunner{},
		&wails_pro.Runner{},
		&wails_dist.Runner{},
		&wails.BundleRunner{},
	) {
		switch r := ss.(type) {
		case goja.Runner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Run(*ctx, args...)
		case wails.Runner:
			ctx := bind.DefaultCoreFactory{}.Create(ctx)
			return r.Run(&Adapter{ctx})
		}
	}

	return fs.ErrInvalid
}

type Adapter struct{ bind.Core }

func listTargets(src string) (out []app.Manifest, err error) {
	s := source.Providers{
		source.OsFs,
	}.GetSource(src)
	if s == nil {
		return nil, fs.ErrNotExist
	}

	for _, ss := range source.CollectT[app.App](s,
		&html.App{},
		&html.Bundle{},
		&html.Project{},
		&js.App{},
		&js.Bundle{},
		&js.Project{},
	) {
		out = append(out, ss.GetMetadata().Manifest)
	}
	if len(out) == 0 {
		err = fs.ErrInvalid
	}
	return
}
