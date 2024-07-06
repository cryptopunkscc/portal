package goja_dev

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/dist"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/target"
	jsEmbed "github.com/cryptopunkscc/portal/target/js/embed"
	"github.com/cryptopunkscc/portal/target/sources"
	"time"
)

type Runner struct {
	distRunner target.Runner[target.DistJs]
}

func NewRunner(newApi target.NewApi, send target.MsgSend) target.Runner[target.ProjectJs] {
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

	// Wait 1 for npm.RunWatch finish initial build otherwise runner can restart on first launch.
	time.Sleep(1 * time.Second)

	return r.distRunner.Run(ctx, projectJs.DistJs())
}
