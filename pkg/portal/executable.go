package portal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const Name = portal.Name

func Executable() string {
	executable, err := os.Executable()
	if err != nil {
		executable = "portal"
	}
	return executable
}

func CmdCtx(ctx context.Context, args ...string) *exec.Cmd {
	e := Executable()
	log.Println("executable.CmdCtx", e, args)
	var c *exec.Cmd
	if ctx != nil {
		c = exec.CommandContext(ctx, Executable(), args...)
	} else {
		c = exec.Command(Executable(), args...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}

func CmdOpenerCtx(ctx context.Context) func(src string, background bool) error {
	run := func(src string) error {
		// resolve apps from given source
		apps, err := ResolveApps(src)
		if len(apps) == 0 {
			return errors.Join(fmt.Errorf("CmdOpenerCtx no apps found in %s", src), err)
		}

		// execute multiple targets as separate processes
		return Spawn(ctx, apps)
	}

	return func(src string, background bool) error {
		if background {
			go run(src)
			return nil
		}
		return run(src)
	}

}

func Spawn(
	ctx context.Context,
	apps Apps,
) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(apps))
	for _, t := range apps {
		go func(app target.App) {
			defer wg.Done()
			if app.Type() == target.Frontend {
				// TODO implement a better way to spawn backends before frontends
				time.Sleep(200 * time.Millisecond)
			}
			if err := CmdCtx(ctx, app.Path(), "--attach").Run(); err != nil {
				return
			}
		}(t)
	}
	wg.Wait()
	return
}

func Attach(
	ctx context.Context,
	bindings runtime.New,
	app target.App,
) (err error) {
	switch app.Type() {

	case target.Backend:
		if err = goja.NewBackend(bindings(app.Type())).RunFs(app.Files()); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()

	case target.Frontend:
		opt := wails.AppOptions(bindings(app.Type()))
		if err = wails.Run(app, opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}

	default:
		return fmt.Errorf("invalid target: %v", app.Path())
	}
	return
}
