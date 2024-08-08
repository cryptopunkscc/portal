package start

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

type Deps interface {
	Port() target.Port
	JoinTarget() target.Dispatch
	DispatchService() target.Dispatch
}

func Inject(deps Deps) Feat {
	return NewFeat(
		deps.Port(),
		deps.DispatchService(),
		deps.JoinTarget(),
	)
}

type Feat struct {
	port         target.Port
	startService target.Dispatch
	runTarget    target.Dispatch
}

func NewFeat(
	port target.Port,
	startService target.Dispatch,
	runTarget target.Dispatch,
) Feat {
	return Feat{
		port:         port,
		runTarget:    runTarget,
		startService: startService,
	}
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

	if err = f.startService(ctx, "s"); err != nil {
		return
	}

	if err = flow.Retry(ctx, 8*time.Second, func(i int, n int, d time.Duration) (err error) {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
		if err = apphost.Init(); err != nil {
			return
		}
		return f.checkService()
	}); err != nil {
		return
	}

	if src != "" {
		return f.runTarget(ctx, src, args...)
	}
	return
}

func (f Feat) checkService() (err error) {
	request := rpc.NewRequest(id.Anyone, f.port.Base)
	if err = rpc.Command(request, "ping"); err == nil {
		return
	}
	return
}
