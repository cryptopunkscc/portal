package dev

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	common "github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(path string, opt *options.App) (err error) {

	front := path
	path = path + "/dist"
	src, err := project.NewModule(path).PortalNodeModule()
	if err != nil {
		return err
	}
	common.SetupOptions(src, opt)
	if opt.Title == "" {
		opt.Title = "development"
	}
	titleSuffix := front
	if exec, err := os.Getwd(); err == nil {
		titleSuffix = exec
	}
	opt.Title = fmt.Sprintf("%s - %s", opt.Title, titleSuffix)

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
