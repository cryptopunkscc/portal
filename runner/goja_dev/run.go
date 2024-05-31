package goja_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/runner/npm"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Runner struct {
	log        plog.Logger
	distRunner *goja_dist.Runner
}

func NewRunner(ctrlPort string, newApi target.NewApi) *Runner {
	distRunner := goja_dist.NewRunner(ctrlPort, newApi)
	return &Runner{distRunner: distRunner}
}

func (r *Runner) Reload() (err error) {
	return r.distRunner.Reload()
}

func (r *Runner) Run(ctx context.Context, project target.ProjectBackend) (err error) {
	r.log = plog.Get(ctx).Type(r).Set(&ctx)
	r.log.Println("staring dev backend", project.Abs())

	go npm.RunWatch(ctx, project.Path())

	return r.distRunner.Run(ctx, project.DistBackend())
}
