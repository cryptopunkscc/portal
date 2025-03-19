package goja_pro

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	jsEmbed "github.com/cryptopunkscc/portal/core/js/embed"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"time"
)

type runner struct {
	distRunner target.ReRunner[target.DistJs]
}

func NewRunner(newRuntime bind.NewRuntime, send target.MsgSend) target.ReRunner[target.ProjectJs] {
	distRunner := goja_dist.NewRunner(newRuntime, send)
	return &runner{distRunner: distRunner}
}

func (r *runner) ReRun() (err error) {
	return r.distRunner.ReRun()
}

func (r *runner) Run(ctx context.Context, projectJs target.ProjectJs, args ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", projectJs.Manifest().Package, projectJs.Abs())
	defer log.Println("exit", projectJs.Manifest().Package, projectJs.Abs())
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}

	libs := target.Any[target.NodeModule](
		target.Skip("node_modules"),
		target.Try(npm2.Resolve),
	).List(
		source.Embed(jsEmbed.PortalLibFS),
	)
	if len(libs) == 0 {
		log.P().Println("libs are empty")
	}

	build := npm_build.Runner(libs...)
	if err = build(ctx, projectJs); err != nil {
		return
	}

	if err = npm.Watch(ctx, projectJs); err != nil {
		return
	}

	// Wait 1 sec for npm.Watch finish initial build otherwise runner can restart on first launch.
	time.Sleep(1 * time.Second)

	return r.distRunner.Run(ctx, projectJs.Dist(), args...)
}
