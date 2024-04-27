package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/backend"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails/dev"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"sync"
)

func Run(
	bindings runtime.New,
	dir string,
) (err error) {
	d, err := runner.New(dir, runner.DevTargets)
	if err != nil {
		return fmt.Errorf("newRunner: %v", err)
	}

	var frontCtxs []context.Context
	appendFrontCtx := func(ctx context.Context) { frontCtxs = append(frontCtxs, ctx) }
	backendEvents := make(chan backend.Event)
	defer close(backendEvents)
	go func() {
		for range backendEvents {
			for _, ctx := range frontCtxs {
				wailsruntime.WindowReload(ctx)
			}
		}
	}()

	wait := sync.WaitGroup{}

	for _, target := range d.Backends {
		wait.Add(1)
		src := ""
		src, err = runner.ResolveSrc(target.Path, "main.js")
		if err != nil {
			return fmt.Errorf("resolveSrc %v: %v", "main.js", err)
		}

		go backend.Watcher(target.Path)

		if err = backend.Dev(goja.NewBackend(context.TODO()), src, backendEvents); err != nil {
			return fmt.Errorf("backend.Dev: %v", err)
		}
	}

	// TODO handle more than one frontend
	for _, target := range d.Frontends {
		wait.Add(1)
		opt := wails.AppOptions(bindings())
		opt.OnStartup = appendFrontCtx
		if err = dev.Run(target.Path, opt); err != nil {
			log.Fatal(fmt.Errorf("dev.Run: %v", err))
		}
		return
	}
	wait.Wait()
	log.Println("dev closed")
	return
}
