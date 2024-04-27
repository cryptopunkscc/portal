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
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"log"
	"os"
	"os/exec"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
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

	return RunTargets(ctx, bindings, targets)
}

func RunTargets(
	ctx context.Context,
	bindings runtime.New,
	targets []runner.Target,
) (err error) {

	// execute single target in current process
	if len(targets) == 1 {
		return RunTarget(ctx, bindings, targets[0])
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
	ctx context.Context,
	bindings runtime.New,
	target runner.Target,
) (err error) {
	switch {

	case runner.IsBackend(target.Files):
		if err = goja.NewBackend(ctx).RunFs(target.Files); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()

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
