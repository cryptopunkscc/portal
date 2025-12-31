package wails_dist

import (
	"context"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/dev"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/target/dev/reload"
)

type Runner struct {
	wails.AppRunner
	send    dev.SendMsg
	newCore bind.NewCore
}

func (r Runner) New() source.Source {
	return &r
}

func (r *Runner) Reload(ctx context.Context) (err error) {
	log := plog.Get(ctx).Type(r)
	if err := r.send(dev.NewMsg(r.Package, dev.Changed)); err != nil {
		log.F().Println(err)
	}
	err = r.Reload(ctx)
	if err := r.send(dev.NewMsg(r.Package, dev.Refreshed)); err != nil {
		log.F().Println(err)
	}
	return err
}

func (r *Runner) Run(ctx bind.Core) (err error) {
	defer plog.TraceErr(&err)
	if !filepath.IsAbs(r.Path) {
		return plog.Errorf("wails_dist.Runner needs absolute path: %s", r.Path)
	}
	plog.Get(ctx).Type(r).Printf("run %T %s", r.App, r.Path)
	go runner.ReloadOnChange(ctx, r, r.Dist)
	r.send = reload.Start(ctx, r.Path, r.Reload, ctx)
	return r.AppRunner.Run(ctx)
}
