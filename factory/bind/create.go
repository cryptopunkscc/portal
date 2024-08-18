package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/target"
)

var ApphostDefault = apphost.Default

func DefaultRuntime() bind.NewRuntime  { return newRuntime(DefaultApphost(ApphostDefault())) }
func FrontendRuntime() bind.NewRuntime { return newRuntime(FrontendApphost(ApphostDefault())) }
func BackendRuntime() bind.NewRuntime  { return newRuntime(BackendApphost(ApphostDefault())) }

func newRuntime(newApphost NewApphost) bind.NewRuntime {
	return func(ctx context.Context, portal target.Portal_) bind.Runtime {
		return bind.Module{
			Apphost: newApphost(ctx, portal),
			Sys:     bind.Sys(ctx),
		}
	}
}
