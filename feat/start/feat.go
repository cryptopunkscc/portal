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
	JoinTarget() target.Request
	Serve() target.Request
}

func Feat(deps Deps) target.Request {
	return feat{
		port:      deps.Port(),
		runTarget: deps.JoinTarget(),
		serve:     deps.Serve(),
	}.Request
}

type feat struct {
	port      target.Port
	serve     target.Request
	runTarget target.Request
}

func (f feat) Request(
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

	if err = f.serve(ctx, "s"); err != nil {
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

func (f feat) checkService() (err error) {
	request := rpc.NewRequest(id.Anyone, f.port.Base)
	if err = rpc.Command(request, "ping"); err == nil {
		return
	}
	return
}
