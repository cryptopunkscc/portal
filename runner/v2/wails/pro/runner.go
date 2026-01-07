package wails_pro

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/v2/wails"
	"github.com/cryptopunkscc/portal/source/html"
	"github.com/cryptopunkscc/portal/target/dev/reload"
	"github.com/wailsapp/wails/v2/pkg/application"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Runner struct {
	html.Project
	frontCtx context.Context
	Core     bind.Core
}

func NewRunner(core bind.Core) (r *Runner) {
	r = &Runner{}
	r.Core = core
	return
}

func (r *Runner) Run(ctx context.Context) (err error) {
	// TODO pass args to js
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", r.Package, r.Path)
	defer log.Println("exit", r.Package, r.Path)

	if err = deps.Check("npm", "-v"); err != nil {
		return
	}

	if err = r.Build(); err != nil {
		return
	}

	opt := wails.AppOptions(r.Core)
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	path := r.Path

	// Start frontend dev watcher
	viteCommand := "npm run dev"
	stopDevWatcher, url, _, err := runViteWatcher(viteCommand, path, true)
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

	// setup opt
	front := path
	path = path + "/dist"
	wails.SetupOptions(opt, r.App)
	if opt.Title == "" {
		opt.Title = "development"
	}
	titleSuffix := front
	if exec, err := os.Getwd(); err == nil {
		titleSuffix = exec
	}
	opt.Title = fmt.Sprintf("%s - %s", opt.Title, titleSuffix)
	opt.LogLevel = 6

	// Setup dev environment
	_ = os.Setenv("devserver", "127.0.0.1:34115")
	_ = os.Setenv("assetdir", path)
	_ = os.Setenv("frontenddevserverurl", url)

	// run
	log.Println("running wails")
	app := application.NewWithOptions(opt)
	go func() {
		<-ctx.Done()
		app.Quit()
	}()

	_ = reload.Start(ctx, r.Package, r.Reload, r.Core)

	err = app.Run()
	if err != nil {
		log.F().Printf("dev.Run: %v", err)
	}
	return
}

func (r *Runner) Reload(ctx context.Context) (err error) {
	if r.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}
