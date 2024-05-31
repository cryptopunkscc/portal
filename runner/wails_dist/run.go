package wails_dist

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/broadcast"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/watcher"
	"github.com/cryptopunkscc/go-astral-js/target"
	"path"
)

type Runner struct {
	ctrlPort string
	inner    *wails.Runner
}

func NewRunner(ctrlPort string, newApi target.NewApi, prefix ...string) (runner *Runner) {
	runner = &Runner{}
	runner.ctrlPort = ctrlPort
	runner.inner = wails.NewRunner(newApi, prefix...)
	return
}

func (r *Runner) Reload() (err error) {
	return r.inner.Reload()
}

func (r *Runner) Run(ctx context.Context, dist target.DistFrontend) (err error) {
	if !path.IsAbs(dist.Abs()) {
		return plog.Errorf("Runner needs absolute path: %s", dist.Abs())
	}

	if err = r.inner.Run(ctx, dist); err != nil {
		return
	}

	pkg := dist.Manifest().Package
	watch := watcher.NewRunner[target.DistFrontend](func() (err error) {
		_ = broadcast.Send(r.ctrlPort, broadcast.NewMsg(pkg, broadcast.Changed))
		err = r.inner.Reload()
		_ = broadcast.Send(r.ctrlPort, broadcast.NewMsg(pkg, broadcast.Refreshed))
		return err
	})

	return watch.Run(ctx, dist)
}
