package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/sig"
	backend "github.com/cryptopunkscc/go-astral-js/pkg/backend/dev"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"os"
	"path"
)

type Backend struct {
	ctx context.Context
	runtime.New
	target.Project
	events sig.Queue[any]
}

func NewBackend(ctx context.Context, bindings runtime.New, project target.Project) *Backend {
	return &Backend{ctx: ctx, New: bindings, Project: project}
}

func (b *Backend) Start() (err error) {
	log.Println("staring dev backend", b.Path())
	src := ""
	if src, err = ResolveSrc(b.Path(), "main.js"); err != nil {
		return fmt.Errorf("resolveSrc %v: %v", "main.js", err)
	}

	go backend.NpmRunWatch(b.ctx, b.Path())
	go b.serve()

	back := goja.NewBackend(b.New(target.Backend, "dev"))
	output := func(event backend.Event) { b.events.Push(event) }
	if err = backend.Dev(b.ctx, back, src, output); err != nil {
		return fmt.Errorf("backend.Dev: %v", err)
	}
	return
}

func (b *Backend) serve() {
	port := devPort(b)
	s := rpc.NewApp(port)
	s.Logger(log.New(log.Writer(), port+" ", 0))
	s.RouteFunc("events", b.events.Subscribe)
	err := s.Run(b.ctx)
	if err != nil {
		log.Printf("%s: %v", port, err)
	}
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
