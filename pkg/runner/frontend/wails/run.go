package wails

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	binding "github.com/cryptopunkscc/go-astral-js/pkg/binding/wails"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io"
	"io/fs"
	"log"
)

func RunFS(src fs.FS, opt *options.App) (err error) {
	SetupOptions(src, opt)
	log.Println("running wails")
	return wails.Run(opt)
}

func SetupOptions(src fs.FS, opt *options.App) {
	// Setup defaults
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	apphostJsFs := binding.WailsJsFs

	// Setup manifest
	log.Println("Reading manifest")
	if manifest, err := bundle.ReadManifestFs(src); err == nil {
		log.Printf("Reading manifest success: %s\n", manifest)
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

func AppOptions(app io.Closer) *options.App {
	return &options.App{
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind:             []interface{}{app},
		OnDomReady: func(ctx context.Context) {
			_ = app.Close()
		},
	}
}
