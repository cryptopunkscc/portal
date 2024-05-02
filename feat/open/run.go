package open

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"os"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
) (err error) {

	apps, err := ResolveApps(src)
	if len(apps) == 0 {
		err = errors.Join(fmt.Errorf("no apps found in %s", src), err)
		return
	}

	// execute single target in current process
	if len(apps) == 1 {
		for _, app := range apps {
			return RunTarget(ctx, bindings, app)
		}
	}

	// execute multiple targets as separate processes
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log.Println("apps", apps)

	for _, t := range apps {
		go func(t target.Source) {
			log.Println("running", t.Path())
			err = portal.Open(ctx, t.Path()).Run()
			_ = exec.Shutdown()
			cancel()
		}(t)
	}

	<-ctx.Done()
	return
}

type Apps map[string]target.App

func ResolveApps(src string) (apps Apps, err error) {
	apps = map[string]target.App{}

	if fs.Exists(src) {
		// scan src as path for portal apps
		apps, err = ResolveAppsByPath(src)
	} else {
		// resolve app path from appstore using given src as package name
		apps[src], err = ResolveAppByPackageName(src)
	}
	return
}

func ResolveAppsByPath(src string) (apps Apps, err error) {
	apps = map[string]target.App{}
	var base, sub string
	base, sub, err = project.Path(src)
	if err != nil {
		return nil, fmt.Errorf("cannot portal apps path: %v", err)
	}
	for app := range project.Find[target.App](os.DirFS(base), sub) {
		apps[app.Manifest().Package] = app
	}
	return
}

func ResolveAppByPackageName(src string) (app target.App, err error) {
	if src, err = appstore.Path(src); err != nil {
		return
	}
	var bundle *project.Bundle
	if bundle, err = project.NewModule(src).Bundle(); err != nil {
		return
	}
	app = bundle
	return
}

func RunTarget(
	ctx context.Context,
	bindings runtime.New,
	app target.App,
) (err error) {
	switch app.Type() {

	case target.Backend:
		if err = goja.NewBackend(ctx).RunFs(app.Files()); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()

	case target.Frontend:
		opt := wails.AppOptions(bindings())
		if err = wails.Run(app, opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}

	default:
		return fmt.Errorf("invalid target: %v", app.Path())
	}
	return
}
