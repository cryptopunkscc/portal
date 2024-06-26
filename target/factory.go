package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"reflect"
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
			plog.Get(ctx).P().Println("cannot create target.NewApi unknown type:", reflect.TypeOf(p))
		}
		a = wrap(n(ctx, p))
		return
	}
}
