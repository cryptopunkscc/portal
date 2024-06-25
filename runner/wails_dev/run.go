package wails_dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/target"
	js "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/project"
	"github.com/cryptopunkscc/portal/target/sources"
	"github.com/wailsapp/wails/v2/pkg/application"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	"syscall"
)

func NewRunner(newApi target.NewApi) *Runner {
	return &Runner{NewApi: newApi}
}

type Runner struct {
	frontCtx context.Context
	target.NewApi
}

func (r *Runner) Run(ctx context.Context, projectHtml target.ProjectHtml) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", projectHtml.Manifest().Package, projectHtml.Abs())
	defer log.Println("exit", projectHtml.Manifest().Package, projectHtml.Abs())

	dependencies := sources.FromFS[target.NodeModule](js.PortalLibFS)
	if err = dist.NewRun(dependencies)(ctx, projectHtml); err != nil {
		return
	}

	api := r.NewApi(ctx, projectHtml)
	opt := wails.AppOptions(api)
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
	src, err := project.FromPath(front)
	if err != nil {
		return err
	}
	wails.SetupOptions(src, opt)
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
	err = app.Run()

	if err != nil {
		log.F().Printf("dev.Run: %v", err)
	}
	return
}

func (r *Runner) Reload() (err error) {
	if r.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(r.frontCtx)
	return
}
