package start

import (
	"context"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	runtime "github.com/cryptopunkscc/portal/runtime/portal"
	"time"
)

type Deps interface {
	Port() target.Port
	Serve() target.Request
	Request() target.Request
}

func Feat(deps Deps) target.Request {
	return feat{
		port:    deps.Port(),
		serve:   deps.Serve(),
		request: deps.Request(),
		portal:  runtime.Client(deps.Port().Base),
	}.Request
}

type feat struct {
	port    target.Port
	serve   target.Request
	request target.Request
	portal  portal.Client
}

func (f feat) Request(
	ctx context.Context,
	src string,
	args ...string,
) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)

	if src == "" {
		err = f.portal.Ping()
	} else {
		err = f.request(ctx, src, args...)
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
		return f.portal.Ping()
	}); err != nil {
		return
	}

	if src != "" {
		return f.request(ctx, src, args...)
	}
	return
}
