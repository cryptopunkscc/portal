package target

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
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
		case Frontend:
			n = frontendApphost
		case Backend:
			n = backendApphost
		default:
			plog.Get(ctx).P().Println("cannot create target.NewApi unknown type:", reflect.TypeOf(p))
		}
		pkg := p.Manifest().Package
		a = wrap(n(ctx, pkg))
		return
	}
}

type NewApphost func(ctx context.Context, pkg string) Apphost
