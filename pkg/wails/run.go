package wails

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/out/wails"
	bindings "github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io/fs"
	"log"
	"os"
)

func Run(src target.App, opt *options.App) (err error) {
	log.Printf("portal open: (%d) %s\n", os.Getpid(), src.Manifest())
	SetupOptions(src, opt)
	app := application.NewWithOptions(opt)
	err = app.Run()
	log.Printf("portal close: (%d) %s\n", os.Getpid(), src.Manifest())
	if err != nil {
		return
	}
	return
}

func SetupOptions(src target.App, opt *options.App) {
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

	apphostJsFs := binding.WailsJsFs

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{Array: []fs.FS{src.Files(), apphostJsFs}}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: src.Files()},
			&assets.FsStore{FS: apphostJsFs}},
		},
	}
}

func AppOptions(app bindings.Api) *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []interface{}{app},
		OnDomReady: func(ctx context.Context) {
			_ = app.Interrupt
		},
	}
}
