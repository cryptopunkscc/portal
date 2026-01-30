package wails

import (
	"context"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/js/embed/wails"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/html"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Runner interface {
	Run(ctx bind.Core) error
}

type BundleRunner struct {
	AppRunner
	Bundle html.Bundle
}

func (r BundleRunner) New() source.Source {
	return &r
}

func (r *BundleRunner) ReadSrc(src source.Source) (err error) {
	if err = r.Bundle.ReadSrc(src); err == nil {
		r.App = r.Bundle.App
		r.Func = r.Run
	}
	return
}

type AppRunner struct {
	html.App
	frontCtx context.Context
}

func (r AppRunner) New() source.Source {
	return &r
}

func (r *AppRunner) ReadSrc(src source.Source) (err error) {
	if err = r.App.ReadSrc(src); err != nil {
		r.Func = r.Run
	}
	return
}

func (r *AppRunner) Reload() (err error) {
	if r.frontCtx == nil {
		return plog.Errorf("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}

func (r *AppRunner) Run(ctx bind.Core) (err error) {
	log := plog.Get(ctx).Type(r)
	log.Println("start", r.Metadata.Package, r.Metadata)
	defer log.Println("exit", r.Metadata.Package, r.Metadata)

	opt := AppOptions(ctx)
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	SetupOptions(opt, r.App)
	if err = application.NewWithOptions(opt).Run(); err != nil {
		return plog.Err(err)
	}
	return
}

func AppOptions(core bind.Core) *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: options.NewRGB(27, 38, 54),
		Bind:             []interface{}{core},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			core.Interrupt()
			return false
		},
	}
}

func SetupOptions(opt *options.App, app html.App) {
	// Setup defaults
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup manifest
	opt.Title = app.Title
	if opt.Title == "" {
		opt.Title = app.Name
	}

	apphostJsFS := wails.JsFs
	assetsFs := app.PathFS()

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{
		assetsFs,
		apphostJsFS,
	}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: assetsFs},
			&assets.FsStore{FS: apphostJsFS},
		}},
	}
}
