package api

import (
	"context"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apphost"
)

type Deps interface {
	TargetDispatch() Dispatch
	WrapApi(Api) Api
}

func New(deps Deps) NewApi {
	apphost.ConnectionsThreshold = 0
	newApphost := apphost.NewFactory(deps.TargetDispatch())
	return apiFactory(deps.WrapApi,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
}

func apiFactory(
	wrap func(Api) Api,
	frontendApphost NewApphost,
	backendApphost NewApphost,
) func(context.Context, Portal_) Api {
	return func(ctx context.Context, p Portal_) (a Api) {
		var n NewApphost
		switch any(p).(type) {
		case PortalHtml:
			n = frontendApphost
		case PortalJs:
			n = backendApphost
		default:
			return
		}
		a = wrap(n(ctx, p))
		return
	}
}
