package wails

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/out/wails"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	bindings "github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io/fs"
	"log"
	"os"
)

func RunFS(src fs.FS, opt *options.App) (err error) {
	m, _ := bundle.ReadManifestFs(src)
	log.Printf("portal open: (%d) %s\n", os.Getpid(), m)
	SetupOptions(src, opt)
	app := application.NewWithOptions(opt)
	err = app.Run()
	log.Printf("portal close: (%d) %s\n", os.Getpid(), m)
	if err != nil {
		return
	}
	return
}

func SetupOptions(src fs.FS, opt *options.App) {
	// Setup defaults
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	apphostJsFs := binding.WailsJsFs

	// Setup manifest
	manifest, err := bundle.ReadManifestFs(src)
	if err == nil {
		opt.Title = manifest.Title
		if opt.Title == "" {
			opt.Title = manifest.Name
		}
	} else {
		log.Println("Reading manifest err: ", err)
	}

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{Array: []fs.FS{src, apphostJsFs}}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: src},
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
