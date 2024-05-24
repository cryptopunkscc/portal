package wails

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/js/embed/wails"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"os"
	"reflect"
)

type Runner struct {
	bindings target.NewApi
	prefix   []string
	log      plog.Logger
}

func NewRunner(newApi target.NewApi, prefix ...string) target.Run[target.AppFrontend] {
	return Runner{bindings: newApi, prefix: prefix}.Run
}

func (r Runner) Run(ctx context.Context, app target.AppFrontend) (err error) {
	r.log = plog.Get(ctx).Type(r).Set(&ctx)
	r.log.Println("Attach frontend", reflect.TypeOf(app), app.Path(), app.Type())
	opt := AppOptions(r.bindings(ctx, app))
	if err = r.run(app, opt); err != nil {
		return fmt.Errorf("dev.Run: %v", err)
	}
	return
}

func (r Runner) run(src target.App, opt *options.App) (err error) {
	r.log.Printf("portal open: (%d) %s", os.Getpid(), src.Manifest())
	defer r.log.Printf("portal close: (%d) %s", os.Getpid(), src.Manifest())
	SetupOptions(src, opt)
	app := application.NewWithOptions(opt)
	err = app.Run()
	return
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
