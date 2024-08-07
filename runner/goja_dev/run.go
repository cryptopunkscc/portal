package goja_dev

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	npm2 "github.com/cryptopunkscc/portal/resolve/npm"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/runner/npm"
	"github.com/cryptopunkscc/portal/runner/npm_build"
	"github.com/cryptopunkscc/portal/target"
	jsEmbed "github.com/cryptopunkscc/portal/target/js/embed"
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

	if err = npm_build.NewRunner(libs...).Run(ctx, projectJs); err != nil {
		return
	}

	if err = npm.RunWatchStart(ctx, projectJs); err != nil {
		return
	}

	// Wait 1 for npm.RunWatchStart finish initial build otherwise runner can restart on first launch.
	time.Sleep(1 * time.Second)

	return r.distRunner.Run(ctx, projectJs.Dist())
}
