package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

var ApphostDefault = apphost.Full

func DefaultRuntime() bind.NewRuntime  { return newRuntime(DefaultApphost(ApphostDefault)) }
func FrontendRuntime() bind.NewRuntime { return newRuntime(FrontendApphost(ApphostDefault)) }
func BackendRuntime() bind.NewRuntime  { return newRuntime(BackendApphost(ApphostDefault)) }

func newRuntime(newApphost NewApphost) bind.NewRuntime {
	return func(ctx context.Context, portal target.Portal_) bind.Runtime {
		return bind.Module{
			Apphost: newApphost(ctx, portal),
			Sys:     bind.Sys(ctx),
		}
	}
}
