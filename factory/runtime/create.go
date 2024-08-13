package runtime

import (
	"context"
	"github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/target"
)

var (
	Default  = newRuntime(apphost.Default())
	Frontend = newRuntime(apphost.Frontend())
	Backend  = newRuntime(apphost.Backend())
)

func newRuntime(n target.NewApphost) target.NewRuntime {
	return func(ctx context.Context, portal target.Portal_) target.Runtime {
		return n(ctx, portal)
	}
}
