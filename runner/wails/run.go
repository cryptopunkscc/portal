package wails

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/js/embed/wails"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type reRunner struct {
	newRuntime bind.NewRuntime
	frontCtx   context.Context
}

func ReRunner(newRuntime bind.NewRuntime) target.ReRunner[target.AppHtml] {
	return &reRunner{newRuntime: newRuntime}
}

func Runner(newRuntime bind.NewRuntime) target.Run[target.AppHtml] {
	return ReRunner(newRuntime).Run
}

func (r *reRunner) ReRun() (err error) {
	if r.frontCtx == nil {
		return plog.Errorf("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}

func (r *reRunner) Run(ctx context.Context, app target.AppHtml, args ...string) (err error) {
	// TODO pass args to js
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", app.Manifest().Package, app.Abs())
	defer log.Println("exit", app.Manifest().Package, app.Abs())
	runtime, ctx := r.newRuntime(ctx, app)
	opt := AppOptions(runtime)
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	SetupOptions(app, opt)
	if err = application.NewWithOptions(opt).Run(); err != nil {
		return plog.Err(err)
	}
	return
}

func AppOptions(runtime bind.Runtime) *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: options.NewRGB(27, 38, 54),
		Bind:             []interface{}{runtime},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			runtime.Interrupt()
			return false
		},
	}
}

func SetupOptions(src target.Portal_, opt *options.App) {
	// Setup defaults
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup manifest
	m := src.Manifest()
	opt.Title = m.Title
	if opt.Title == "" {
		opt.Title = m.Name
	}

	apphostJsFs := wails.JsFs

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{src.Files(), apphostJsFs}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: src.Files()},
			&assets.FsStore{FS: apphostJsFs}},
		},
	}
}
