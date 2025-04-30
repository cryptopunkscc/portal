package goja_pro

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/deps"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
	"github.com/cryptopunkscc/portal/target/js"
	"github.com/cryptopunkscc/portal/target/npm"
	"time"
)

func Runner(newCore bind.NewCore) *target.SourceRunner[target.ProjectJs] {
	return &target.SourceRunner[target.ProjectJs]{
		Resolve: target.Any[target.ProjectJs](js.ResolveProject.Try),
		Runner: &ReRunner{
			distRunner: &goja_dist.ReRunner{
				NewCore: newCore,
			},
		},
	}
}

type ReRunner struct {
	distRunner target.ReRunner[target.DistJs]
}

func (r *ReRunner) Reload() (err error) {
	return r.distRunner.Reload()
}

func (r *ReRunner) Run(ctx context.Context, projectJs target.ProjectJs, args ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Println("start", projectJs.Manifest().Package, projectJs.Abs())
	defer log.Println("exit", projectJs.Manifest().Package, projectJs.Abs())
	if err = deps.RequireBinary("npm"); err != nil {
		return
	}

	if err = npm.BuildProject().Run(ctx, projectJs); err != nil {
		return
	}

	if err = npm.Watch(ctx, projectJs); err != nil {
		return
	}

	// Wait 1 sec for npm.Watch finish initial build otherwise runner can restart on first launch.
	time.Sleep(1 * time.Second)

	return r.distRunner.Run(ctx, projectJs.Dist(), args...)
}
