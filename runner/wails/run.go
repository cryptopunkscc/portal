package wails

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/js/embed/wails"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Runner struct {
	newApi   target.NewApi
	frontCtx context.Context
}

func NewRunner(newApi target.NewApi) *Runner {
	return &Runner{newApi: newApi}
}

func (r *Runner) Reload() (err error) {
	if r.frontCtx == nil {
		return plog.Errorf("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}

func (r *Runner) Run(ctx context.Context, app target.AppHtml) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", app.Manifest().Package, app.Abs())
	defer log.Println("exit", app.Manifest().Package, app.Abs())
	api := r.newApi(ctx, app)
	opt := AppOptions(api)
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	SetupOptions(app, opt)
	if err = application.NewWithOptions(opt).Run(); err != nil {
		return plog.Err(err)
	}
	return
}

func AppOptions(app target.Api) *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []interface{}{app},
		OnDomReady: func(ctx context.Context) {
			app.Interrupt()
		},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			app.Interrupt()
			return false
		},
	}
}

func SetupOptions(src target.Portal, opt *options.App) {
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
