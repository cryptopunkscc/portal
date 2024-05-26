package dispatch

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	osexec "os/exec"
	"time"
)

type Feat struct {
	dispatch   target.Dispatch
	executable string
}

func NewFeat(executable string, dispatch target.Dispatch) target.Dispatch {
	return Feat{
		dispatch:   dispatch,
		executable: executable,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	src string,
	args ...string,
) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)

	if err = f.dispatch(ctx, src, args...); err == nil {
		return
	}

	if err = f.portalServe(ctx); err != nil {
		return
	}

	if err = exec.Retry(ctx, 8*time.Second, func(i int, n int, duration time.Duration) error {
		return f.dispatch(ctx, src, args...)
	}); err != nil {
		return
	}

	return
}

func (f Feat) portalServe(ctx context.Context) (err error) {
	c := osexec.CommandContext(ctx, f.executable, "s", "-t")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
