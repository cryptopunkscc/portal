package wails_pro

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target/html"
	"github.com/wailsapp/wails/v2/pkg/application"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	"syscall"
)

func Runner(newCore bind.NewCore) *target.SourceRunner[target.ProjectHtml] {
	return &target.SourceRunner[target.ProjectHtml]{
		Resolve: target.Any[target.ProjectHtml](html.ResolveProject.Try),
		Runner:  ReRunner(newCore),
	}
}

func ReRunner(newCore bind.NewCore) target.ReRunner[target.ProjectHtml] {
	return &reRunner{NewCore: newCore}
}

type reRunner struct {
	frontCtx context.Context
	bind.NewCore
}

func (r *reRunner) Run(ctx context.Context, projectHtml target.ProjectHtml, args ...string) (err error) {
	// TODO pass args to js
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", projectHtml.Manifest().Package, projectHtml.Abs())
	defer log.Println("exit", projectHtml.Manifest().Package, projectHtml.Abs())

	if err = deps.RequireBinary("npm"); err != nil {
		return
	}

	build := npm_build.NewRun()
	if err = build(ctx, projectHtml); err != nil {
		return
	}

	core, ctx := r.NewCore(ctx, projectHtml)
	opt := wails.AppOptions(core)
	opt.OnStartup = func(ctx context.Context) { r.frontCtx = ctx }
	path := projectHtml.Abs()

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
	wails.SetupOptions(projectHtml, opt)
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
	_ = os.Setenv("devserver", "localhost:34115")
	_ = os.Setenv("assetdir", path)
	_ = os.Setenv("frontenddevserverurl", url)

	// run
	log.Println("running wails")
	app := application.NewWithOptions(opt)
	go func() {
		<-ctx.Done()
		app.Quit()
	}()

	_ = reload.Start(ctx, projectHtml, r.Reload, core)

	err = app.Run()
	if err != nil {
		log.F().Printf("dev.Run: %v", err)
	}
	return
}

func (r *reRunner) Reload() (err error) {
	if r.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}
