package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/target/dev/reload"
	"github.com/cryptopunkscc/portal/target/dist"
	"github.com/cryptopunkscc/portal/target/js"
	"path/filepath"
	"time"
)

func Runner(newCore bind.NewCore) *target.SourceRunner[target.DistJs] {
	return &target.SourceRunner[target.DistJs]{
		Resolve: target.Any[target.DistJs](
			js.ResolveDist.Try,
			js.ResolveBundle.Try,
		),
		Runner: &ReRunner{NewCore: newCore},
	}
}

type ReRunner struct {
	bind.NewCore
	send    target.MsgSend
	dist    target.DistJs
	backend *goja.Backend
}

func (r *ReRunner) Reload() (err error) {
	return r.backend.RunFs(r.dist.FS())
}

func (r *ReRunner) Run(ctx context.Context, distJs target.DistJs, args ...string) (err error) {
	if any(r.NewCore) == nil {
		panic("newCore cannot be nil")
	}
	if !filepath.IsAbs(distJs.Abs()) {
		return plog.Errorf("ReRunner needs absolute path: %s", distJs.Abs())
	}
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("run %T %s", distJs, distJs.Abs())
	core, ctx := r.NewCore(ctx, distJs)
	r.backend = goja.NewBackend(core)
	r.dist = distJs
	if err = r.Reload(); err != nil {
		log.E().Println(err.Error())
	}
	pkg := distJs.Manifest().Package
	watch := dist.ReRunner[target.DistJs](func(...string) error {
		if err := r.send(target.NewMsg(pkg, target.DevChanged)); err != nil {
			log.E().Println(err)
		}
		if err := r.Reload(); err != nil {
			log.E().Println(err.Error())
		}
		// TODO find better solution then sleep
		// target.DevRefreshed msg must be delayed until backend is fully refreshed (all ports registered).
		time.Sleep(2 * time.Second)
		if err := r.send(target.NewMsg(pkg, target.DevRefreshed)); err != nil {
			log.E().Println(err)
		}
		return nil
	})
	r.send = reload.Start(ctx, distJs, r.Reload, core)
	return watch.Run(ctx, distJs, args...)
}
