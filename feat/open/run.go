package open

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"log"
	"os"
	"os/exec"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
) (err error) {
	var targets []runner.Target
	if fs.Exists(src) {
		for target := range project.ProdTargets(os.DirFS(src)) {
			targets = append(targets, target)
		}
	} else {
		if src, err = appstore.Path(src); err != nil {
			return
		}
		var bundle *project.Bundle
		if bundle, err = project.NewModule(src).Bundle(); err != nil {
			return
		}
		targets = append(targets, bundle)
	}

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
	switch target.Type() {

	case runner.Backend:
		if err = goja.NewBackend(ctx).RunFs(target.Files()); err != nil {
			return fmt.Errorf("goja.NewBackend().RunSource: %v", err)
		}
		<-ctx.Done()

	case runner.Frontend:
		opt := wails.AppOptions(bindings())
		if err = wails.RunFS(target.Files(), opt); err != nil {
			return fmt.Errorf("dev.Run: %v", err)
		}

	default:
		return fmt.Errorf("invalid target: %v", target.Path())
	}
	return
}

func RunTargetProcess(ctx context.Context, target runner.Target) (err error) {
	log.Println("RunTargetProcess: ", target.Path())
	cmd := exec.CommandContext(ctx, portal.Executable(), target.Path())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
