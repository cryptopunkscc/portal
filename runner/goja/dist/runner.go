package goja_dist

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/cryptopunkscc/portal/api/dev"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/target/dev/reload"
)

type Runner struct {
	goja.AppRunner
	send dev.SendMsg
}

func (r Runner) New() source.Source {
	return &r
}

func (r *Runner) Reload(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(r)
	if err := r.send(dev.NewMsg(r.Package, dev.Changed)); err != nil {
		log.E().Println(err)
	}
	err = r.AppRunner.Reload(ctx)
	time.Sleep(2 * time.Second) // target.DevRefreshed msg must be delayed until backend is fully refreshed (all ports registered). TODO find better solution then sleep
	if err := r.send(dev.NewMsg(r.Package, dev.Refreshed)); err != nil {
		log.E().Println(err)
	}
	return
}

func (r *Runner) Run(ctx bind.Context, args ...string) (err error) {
	defer plog.TraceErr(&err)
	if !filepath.IsAbs(r.Path) {
		return fmt.Errorf("goja_dist.Runner needs absolute path: %s", r.Path)
	}
	log := plog.Get(ctx).Type(r)
	log.Printf("run %T %s", r.App, r.Path)
	r.Args = args
	if err = r.AppRunner.Reload(ctx); err != nil {
		log.E().Println(err.Error())
	}
	r.send = reload.Start(ctx, r.Package, r.Reload, r.Core)
	return runner.ReloadOnChange(ctx, r, r.Dist)
}
