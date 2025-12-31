package wails

import (
	"context"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/js/embed/wails"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Runner interface {
	Run(ctx context.Context) error
}

type BundleRunner struct {
	AppRunner
	Bundle source.HtmlBundle
}

func NewBundleRunner(core bind.Core) *BundleRunner {
	return &BundleRunner{AppRunner: AppRunner{Core: core}}
}

func (r *BundleRunner) ReadSrc(src source.Source) (err error) {
	if err = r.Bundle.ReadSrc(src); err != nil {
		r.HtmlApp = r.Bundle.HtmlApp
		r.Func = r.Run
	}
	return
}

type AppRunner struct {
	source.HtmlApp
	Core     bind.Core
	frontCtx context.Context
}

func NewAppRunner(core bind.Core) *AppRunner {
	return &AppRunner{Core: core}
}

func (r *AppRunner) ReadSrc(src source.Source) (err error) {
	if err = r.Html.ReadSrc(src); err != nil {
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

func (r *AppRunner) Run(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", r.Metadata.Package, r.Metadata)
	defer log.Println("exit", r.Metadata.Package, r.Metadata)

	opt := r.AppOptions()
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	r.SetupOptions(opt)
	if err = application.NewWithOptions(opt).Run(); err != nil {
		return plog.Err(err)
	}
	return
}

func (r *AppRunner) AppOptions() *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: options.NewRGB(27, 38, 54),
		Bind:             []interface{}{r.Core},
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			r.Core.Interrupt()
			return false
		},
	}
}

func (r *AppRunner) SetupOptions(opt *options.App) {
	// Setup defaults
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup manifest
	opt.Title = r.Metadata.Title
	if opt.Title == "" {
		opt.Title = r.Metadata.Name
	}

	apphostJsFS := wails.JsFs

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{
		r.FS(),
		apphostJsFS,
	}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: r.FS()},
			&assets.FsStore{FS: apphostJsFS}},
		},
	}
}
