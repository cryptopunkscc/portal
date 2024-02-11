package dev

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func Run(path string, opt *options.App) (err error) {

	front := path
	path = path + "/dist"

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

	// Start frontend dev watcher
	runDevWatcherCommand := "npm run dev"
	stopDevWatcher, url, _, err := runFrontendDevWatcherCommand(front, runDevWatcherCommand, true)
	if err != nil {
		return err
	}
	log.Println("url: ", url)
	go func() {
		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, os.Interrupt, syscall.SIGTERM)
		<-quitChannel
		stopDevWatcher()
	}()

	// Setup dev environment
	os.Setenv("devserver", "localhost:34115")
	os.Setenv("assetdir", path)
	os.Setenv("frontenddevserverurl", url)

	// run
	log.Println("running wails")
	return wails.Run(opt)
}
