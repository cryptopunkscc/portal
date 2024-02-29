package dev

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/assets"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(path string, opt *options.App) (err error) {

	front := path
	path = path + "/dist"

	// Setup defaults
	//if opt.Title != "" {
	//	opt.Title = filepath.Base(front)
	//}
	if opt.AssetServer == nil {
		opt.AssetServer = &assetserver.Options{}
	}

	// Setup fs assets
	opt.AssetServer.Assets = assets.ArrayFs{Array: []fs.FS{os.DirFS(path), apphost.JsWailsFs()}}

	// Setup http assets
	opt.AssetServer.Handler = assets.StoreHandler{
		Store: &assets.OverlayStore{Stores: []assets.Store{
			&assets.FsStore{FS: os.DirFS(path)},
			&assets.FsStore{FS: apphost.JsWailsFs()}},
		},
	}

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
	opt.LogLevel = 6
	return wails.Run(opt)
}
