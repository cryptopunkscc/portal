package target

import (
	"context"
)

// ApiFactory deprecated FIXME refactor & remove
func ApiFactory(
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
