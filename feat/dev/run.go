package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/backend"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails"
	"github.com/cryptopunkscc/go-astral-js/pkg/wails/dev"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path"
	"sync"
)

func Run(
	bindings runtime.New,
	dir string,
) (err error) {
	var backends []project.PortalNodeModule
	var frontends []project.PortalNodeModule
	for module := range project.Find[project.PortalNodeModule](os.DirFS(dir), ".") {
		switch module.Type() {
		case target.Frontend:
			frontends = append(frontends, module)
		case target.Backend:
			backends = append(backends, module)
		case target.Invalid:
		}
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

	for _, t := range backends {
		wait.Add(1)
		src := ""
		src, err = ResolveSrc(t.Path(), "main.js")
		if err != nil {
			return fmt.Errorf("resolveSrc %v: %v", "main.js", err)
		}

		go backend.Watcher(t.Path())

		if err = backend.Dev(goja.NewBackend(context.TODO()), src, backendEvents); err != nil {
			return fmt.Errorf("backend.Dev: %v", err)
		}
	}

	// TODO handle more than one frontend
	for _, target := range frontends {
		wait.Add(1)
		opt := wails.AppOptions(bindings())
		opt.OnStartup = appendFrontCtx
		if err = dev.Run(target.Path(), opt); err != nil {
			log.Fatal(fmt.Errorf("dev.Run: %v", err))
		}
		return
	}
	wait.Wait()
	log.Println("dev closed")
	return
}

func ResolveSrc(dir string, name string) (f string, err error) {
	f = path.Join(dir, "dist", name)
	if _, err = os.Stat(f); err == nil {
		return
	}
	f = path.Join(dir, name)
	if _, err = os.Stat(f); err == nil {
		return
	}
	return
}
