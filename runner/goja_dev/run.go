package goja_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/dist"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dist"
	"github.com/cryptopunkscc/go-astral-js/runner/npm"
	"github.com/cryptopunkscc/go-astral-js/target"
	jsEmbed "github.com/cryptopunkscc/go-astral-js/target/js/embed"
	"github.com/cryptopunkscc/go-astral-js/target/sources"
)

type Runner struct {
	distRunner *goja_dist.Runner
}

func NewRunner(newApi target.NewApi, send target.MsgSend) *Runner {
	distRunner := goja_dist.NewRunner(newApi, send)
	return &Runner{distRunner: distRunner}
}

func (r *Runner) Reload() (err error) {
	return r.distRunner.Reload()
}

func (r *Runner) Run(ctx context.Context, projectJs target.ProjectJs) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", projectJs.Manifest().Package, projectJs.Abs())
	defer log.Println("exit", projectJs.Manifest().Package, projectJs.Abs())

	dependencies := sources.FromFS[target.NodeModule](jsEmbed.PortalLibFS)
	if err = dist.NewRun(dependencies)(ctx, projectJs); err != nil {
		return
	}

	if err = npm.RunWatch(ctx, projectJs.Abs()).Start(); err != nil {
		return
	}

	return r.distRunner.Run(ctx, projectJs.DistJs())
}
