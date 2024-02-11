package wails

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"path/filepath"
)

func Run(path string, opt *options.App) (err error) {
	// identify app bundle type
	bundleType, err := assets.BundleType(path)
	if err != nil {
		return
	}

	// Setup defaults
	if opt.Title != "" {
		opt.Title = filepath.Base(path)
	}
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup fs assets
	opt.AssetServer.Assets, err = assets.BundleFS(bundleType, path)
	if err != nil {
		return
	}

	// Setup http assets
	store, err := assets.BundleStore(bundleType, path)
	if err != nil {
		return
	}
	opt.AssetServer.Handler = assets.StoreHandler{Store: store}

	log.Println("running wails")
	return wails.Run(opt)
}
