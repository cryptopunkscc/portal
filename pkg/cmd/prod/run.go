package prod

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner/backend/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner/frontend/wails"
	"log"
	"os"
	"os/exec"
	"sync"
)

func RunLegacy(
	bindings runner.Bindings,
	src string,
) (err error) {
	d, err := runner.New(src, runner.ProdTargets)
	if err != nil {
		return fmt.Errorf("newRunner: %v", err)
	}
	wait := sync.WaitGroup{}

	for _, target := range d.Backends {
		wait.Add(1)
		if err = goja.NewBackend().RunFs(target.Files); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
	}

	// TODO handle more than one frontend
	for _, target := range d.Frontends {
		wait.Add(1)
		opt := wails.AppOptions(bindings())
		if err = wails.RunFS(target.Files, opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}
		return
	}
	wait.Wait()
	return
}

func Run(
	bindings runner.Bindings,
	src string,
) (err error) {

	if !fs.Exists(src) {
		log.Println("resolving path from id: ", src)
		if src, err = appstore.Path(src); err != nil {
			return
		}
	}

	targets, err := runner.ProdTargets(src)
	if err != nil {
		return fmt.Errorf("prod.Run: %v", err)
	}

	return RunTargets(bindings, targets)
}

func RunTargets(
	bindings runner.Bindings,
	targets []runner.Target,
) (err error) {

	// execute single target in current process
	if len(targets) == 1 {
		return RunTarget(bindings, targets[0])
	}

	// execute multiple targets as separate processes
	ctx, cancel := context.WithCancel(context.Background())
	for _, target := range targets {
		go func(target runner.Target) {
			err = RunTargetProcess(ctx, target)
			cancel()
		}(target)
	}
	<-ctx.Done()
	cancel()
	return
}

func RunTarget(
	bindings runner.Bindings,
	target runner.Target,
) (err error) {
	switch {

	case runner.IsBackend(target.Files):
		log.Println("Running in backend mode")
		if err = goja.NewBackend().RunFs(target.Files); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-context.Background().Done()

	case runner.IsFrontend(target.Files):
		opt := wails.AppOptions(bindings())
		if err = wails.RunFS(target.Files, opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}

	default:
		return fmt.Errorf("invalid target: %v", target.Path)
	}
	return
}

func RunTargetProcess(ctx context.Context, target runner.Target) (err error) {
	log.Println("RunTargetProcess: ", target.Path)
	cmd := exec.CommandContext(ctx, portal.Executable(), target.Path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
