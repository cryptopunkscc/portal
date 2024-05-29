package goja_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/backend_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func NewRun(newApi target.NewApi) target.Run[target.ProjectBackend] {
	distRunner := goja_dist.NewRunner(newApi)
	return (&Runner{distRunner: distRunner}).Run
}

type Runner struct {
	log        plog.Logger
	distRunner *goja_dist.Runner
}

func (b *Runner) Run(ctx context.Context, project target.ProjectBackend) (err error) {
	b.log = plog.Get(ctx).Type(b).Set(&ctx)
	b.log.Println("staring dev backend", project.Abs())

	go backend_dev.NpmRunWatch(ctx, project.Path())
	go b.runDevRpc(ctx, project)

	return b.distRunner.Run(ctx, project.DistBackend())
}

func (b *Runner) runDevRpc(ctx context.Context, project target.ProjectBackend) {
	port := target.DevPort(project)
	s := rpc.NewApp(port)
	s.RouteFunc("events", b.distRunner.Events().Subscribe)
	err := s.Run(ctx)
	if err != nil {
		b.log.Printf("%s: %v", port, err)
	}
}
