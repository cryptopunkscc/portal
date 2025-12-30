package bind

import (
	"context"
	"os"

	apphost2 "github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/api/target"
	core "github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

var coreFactory = CoreFactory{}

var NewDefaultCore = coreFactory.NewDefaultFunc()
var NewFrontendCore = coreFactory.NewFrontendFunc()
var NewBackendCore = coreFactory.NewBackendFunc()

type CoreFactory struct{ token.Repository }

func (f CoreFactory) NewDefaultFunc() NewCore  { return f.NewFrontendFunc() }
func (f CoreFactory) NewFrontendFunc() NewCore { return f.api(f.frontendApphost(f.cachedInvoker)) }
func (f CoreFactory) NewBackendFunc() NewCore  { return f.api(f.backendApphost(f.cachedInvoker)) }

func (f CoreFactory) api(newApphost newApphost) NewCore {
	return func(ctx context.Context, portal target.Portal_) (Core, context.Context) {
		m := struct {
			Apphost
			bind.Process
		}{}
		m.Process, ctx = Process(ctx)
		m.Apphost, ctx = newApphost(ctx, portal)
		return m, ctx
	}
}

func (f CoreFactory) cachedInvoker(ctx context.Context, portal target.Portal_) apphost.Cached {
	i := &core.Invoker{Ctx: ctx}
	i.Log = plog.Get(ctx).Type(i)
	if f.Adapter != nil {
		i.Endpoint = f.Adapter.Endpoint
	}

	// TODO deprecated - remove after fixing auth token injection for mobile
	if len(os.Getenv(apphost2.AuthTokenEnv)) == 0 && f.Repository.Load() {
		if t, err := f.Repository.Get(portal.Manifest().Package); err == nil {
			i.AuthToken = string(t.Token)
		}
	}
	return core.Cached(i)
}

type newApphost func(ctx context.Context, portal target.Portal_) (Apphost, context.Context)

func (f CoreFactory) frontendApphost(
	create func(ctx context.Context, portal target.Portal_) apphost.Cached,
) newApphost {
	return func(ctx context.Context, portal target.Portal_) (Apphost, context.Context) {
		return Adapter(ctx, create(ctx, portal), portal.Manifest().Package), ctx
	}
}

func (f CoreFactory) backendApphost(
	create func(ctx context.Context, portal target.Portal_) apphost.Cached,
) newApphost {
	return func(ctx context.Context, portal target.Portal_) (Apphost, context.Context) {
		cached := create(ctx, portal)
		core.ConnectionsThreshold = 0
		ctx = core.Timeout(ctx, cached, portal)
		return Adapter(ctx, cached, portal.Manifest().Package), ctx
	}
}
