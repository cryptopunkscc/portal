package api

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apphost"
)

type Deps interface {
	TargetDispatch() target.Dispatch
	WrapApi(target.Api) target.Api
}

func New(deps Deps) target.NewApi {
	apphost.ConnectionsThreshold = 0
	newApphost := apphost.NewFactory(deps.TargetDispatch())
	return target.ApiFactory(deps.WrapApi,
		newApphost.NewAdapter,
		newApphost.WithTimeout,
	)
}
