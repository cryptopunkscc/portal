package goja_dist

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/runner/backend_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Runner struct {
	target.NewApi
	events sig.Queue[any]
}

func NewRun(newApi target.NewApi) target.Run[target.DistBackend] {
	return NewRunner(newApi).Run
}

func NewRunner(newApi target.NewApi) *Runner {
	return &Runner{NewApi: newApi}
}

func (b *Runner) Events() *sig.Queue[any] {
	return &b.events
}

func (b *Runner) Run(ctx context.Context, dist target.DistBackend) (err error) {
	back := goja.NewBackend(b.NewApi(ctx, dist))
	output := func(event backend_dev.Event) { b.events.Push(event) }
	if err = backend_dev.Dev(ctx, back, dist, output); err != nil {
		return
	}
	return
}
