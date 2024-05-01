package open

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"os"
	"os/exec"
	"path"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
) (err error) {
	var apps []target.App
	log.Println("opening apps...", src)
	if fs.Exists(src) {
		for app := range project.Apps(os.DirFS(src)) {
			apps = append(apps, app)
		}
	} else {
		if src, err = appstore.Path(src); err != nil {
			return
		}
		var bundle *project.Bundle
		if bundle, err = project.NewModule(src).Bundle(); err != nil {
			return
		}
		apps = append(apps, bundle)
	}

	// execute single target in current process
	if len(apps) == 1 {
		return RunTarget(ctx, bindings, apps[0])
	}

	// execute multiple targets as separate processes
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("apps", apps)
	for _, t := range apps {
		go func(t target.Source) {
			err = RunTargetProcess(ctx, path.Join(src, t.Path()))
			cancel()
		}(t)
	}
	<-ctx.Done()
	cancel()
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

func RunTargetProcess(ctx context.Context, src string) (err error) {
	log.Println("RunTargetProcess: ", src)
	cmd := exec.CommandContext(ctx, portal.Executable(), src)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
