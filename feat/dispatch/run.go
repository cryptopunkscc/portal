package dispatch

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"time"
)

type Feat struct {
	dispatchTarget  target.Dispatch
	dispatchService target.Dispatch
}

func NewFeat(
	dispatchTarget target.Dispatch,
	dispatchService target.Dispatch,
) target.Dispatch {
	return Feat{
		dispatchTarget:  dispatchTarget,
		dispatchService: dispatchService,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	src string,
	args ...string,
) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)

	if err = f.dispatchTarget(ctx, src, args...); err == nil {
		return
	}

	if err = f.dispatchService(ctx, "s", "-t"); err != nil {
		return
	}

	if err = exec.Retry(ctx, 8*time.Second, func(i int, n int, duration time.Duration) (err error) {
		if err = apphost.Init(); err != nil {
			return
		}
		return f.dispatchTarget(ctx, src, args...)
	}); err != nil {
		return
	}

	return
}
