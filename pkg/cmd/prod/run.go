package prod

import (
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner/backend/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner/frontend/wails"
	"sync"
)

func Run(
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
