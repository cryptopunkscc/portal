package dispatch

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
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
	log := plog.Get(ctx).Type(f).Set(&ctx)

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

	if err = flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
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
