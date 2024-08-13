package runtime

import (
	"context"
	"github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/target"
)

var (
	Default  = runtime(apphost.Default())
	Frontend = runtime(apphost.Frontend())
	Backend  = runtime(apphost.Backend())
)

func runtime(n target.NewApphost) target.NewApi {
	return func(ctx context.Context, portal target.Portal_) target.Api {
		return n(ctx, portal)
	}
}
