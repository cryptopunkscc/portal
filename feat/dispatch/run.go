package dispatch

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
	"time"
)

type Feat struct {
	port         target.Port
	runTarget    target.Dispatch
	startService target.Dispatch
}

func NewFeat(
	port target.Port,
	runTarget target.Dispatch,
	startService target.Dispatch,
) target.Dispatch {
	return Feat{
		port:         port,
		runTarget:    runTarget,
		startService: startService,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	src string,
	args ...string,
) (err error) {
	plog.Get(ctx).Type(f).Set(&ctx)

	if src == "" {
		err = f.checkService()
	} else {
		err = f.runTarget(ctx, src, args...)
	}

	if err == nil || ctx.Err() != nil {
		return
	}

	if err = f.startService(ctx, "s", "-t"); err != nil {
		return
	}

	if src == "" {
		return
	}

	if err = exec.Retry(ctx, 8*time.Second, func(i int, n int, duration time.Duration) (err error) {
		if err = apphost.Init(); err != nil {
			return
		}
		return f.runTarget(ctx, src, args...)
	}); err != nil {
		return
	}

	return
}

func (f Feat) checkService() (err error) {
	request := rpc.NewRequest(id.Anyone, f.port.String())
	if err = rpc.Command(request, "ping"); err == nil {
		return
	}
	return
}
