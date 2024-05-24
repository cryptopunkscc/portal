package goja_dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/backend_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	"path"
)

type Runner struct {
	target.NewApi
	events sig.Queue[any]
	log    plog.Logger
}

func NewRunner(newApi target.NewApi) target.Run[target.ProjectBackend] {
	return (&Runner{NewApi: newApi}).Run
}

func (b *Runner) Run(ctx context.Context, project target.ProjectBackend) (err error) {
	b.log = plog.Get(ctx).Type(b).Set(&ctx)
	b.log.Println("staring dev backend", project.Abs())
	src := ""
	if src, err = ResolveSrc(project.Path(), "main.js"); err != nil {
		return fmt.Errorf("resolveSrc %v: %v", "main.js", err)
	}

	go backend_dev.NpmRunWatch(ctx, project.Path())
	go b.serve(ctx, project)

	back := goja.NewBackend(b.NewApi(ctx, project))
	output := func(event backend_dev.Event) { b.events.Push(event) }
	if err = backend_dev.Dev(ctx, back, src, output); err != nil {
		return fmt.Errorf("backend.Dev: %v", err)
	}
	return
}

func (b *Runner) serve(ctx context.Context, project target.ProjectBackend) {
	port := target.DevPort(project)
	s := rpc.NewApp(port)
	//s.Logger(log.New(log.Writer(), port+" ", 0))
	s.RouteFunc("events", b.events.Subscribe)
	err := s.Run(ctx)
	if err != nil {
		b.log.Printf("%s: %v", port, err)
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
