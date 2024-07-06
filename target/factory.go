package target

import (
	"context"
)

func ApiFactory(
	wrap func(Api) Api,
	frontendApphost NewApphost,
	backendApphost NewApphost,
) func(context.Context, Portal) Api {
	return func(ctx context.Context, p Portal) (a Api) {
		var n NewApphost
		switch any(p).(type) {
		case Html:
			n = frontendApphost
		case Js:
			n = backendApphost
		default:
			return
		}
		a = wrap(n(ctx, p))
		return
	}
}
