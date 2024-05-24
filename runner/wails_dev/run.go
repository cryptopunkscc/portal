package wails_dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/project"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	"syscall"
)

type Runner struct {
	frontCtx context.Context
	target.NewApi
	log plog.Logger
}

func NewRunner(newApi target.NewApi) target.Run[target.ProjectFrontend] {
	return Runner{NewApi: newApi}.Run
}

func (f Runner) Run(ctx context.Context, project target.ProjectFrontend) (err error) {
	f.log = plog.Get(ctx).Type(f).Set(&ctx)
	f.log.Printf("portal dev open: (%d) %s\n", os.Getpid(), project.Manifest())
	defer f.log.Printf("portal dev close: (%d) %s\n", os.Getpid(), project.Manifest())
	opt := wails.AppOptions(f.NewApi(ctx, project))
	opt.OnStartup = func(ctx context.Context) {
		f.frontCtx = ctx
		go f.serve(project)
	}
	if err = f.run(project.Abs(), opt); err != nil {
		f.log.F().Printf("dev.Run: %v", err)
	}
	return
}

func (f Runner) serve(project target.ProjectFrontend) {
	port := target.DevPort(project)
	s := rpc.NewApp(port)
	//s.Logger(f.log.New(f.log.Writer(), port+" ", 0))
	s.RouteFunc("reload", f.Reload)
	err := s.Run(f.frontCtx)
	if err != nil {
		f.log.Printf("%s: %v", port, err)
	}
}

func (f Runner) Reload() (err error) {
	if f.frontCtx == nil {
		return errors.New("nil context")
	}
	wailsruntime.WindowReload(f.frontCtx)
	return
}

func (f Runner) run(path string, opt *options.App) (err error) {
	// Start frontend dev watcher
	viteCommand := "npm run dev"
	stopDevWatcher, url, _, err := runViteWatcher(viteCommand, path, true)
	if err != nil {
		return err
	}
	f.log.Println("url: ", url)
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
	f.log.Println("running wails")
	app := application.NewWithOptions(opt)
	err = app.Run()
	return
}
