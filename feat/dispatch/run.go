package dispatch

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"time"
)

type Feat struct {
	runTarget  target.Dispatch
	runService target.Dispatch
}

func NewFeat(
	runTarget target.Dispatch,
	runService target.Dispatch,
) target.Dispatch {
	return Feat{
		runTarget:  runTarget,
		runService: runService,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	src string,
	args ...string,
) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)

	if err = f.runTarget(ctx, src, args...); err == nil {
		return
	}

	if err = f.runService(ctx, "s", "-t"); err != nil {
		return
	}

	if err = exec.Retry(ctx, 8*time.Second, func(i int, n int, duration time.Duration) error {
		return f.runTarget(ctx, src, args...)
	}); err != nil {
		return
	}

	return
}
