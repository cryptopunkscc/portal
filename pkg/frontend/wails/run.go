package wails

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io"
	"io/fs"
	"log"
)

func RunFS(src fs.FS, opt *options.App) (err error) {
	// Setup defaults
	//if opt.Title != "" {
	//	opt.Title = filepath.Base(path)
	//}
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{Array: []fs.FS{src, apphost.JsWailsFs()}}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: src},
			&assets.FsStore{FS: apphost.JsWailsFs()}},
		},
	}

	log.Println("running wails")
	return wails.Run(opt)
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
