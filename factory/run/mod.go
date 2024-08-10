package run

import (
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/runtime/api"
	. "github.com/cryptopunkscc/portal/target"
)

type Module[T Portal_] struct {
	Deps[T]
	di.Cache
}

type Deps[T Portal_] interface {
	api.Deps
	open.Deps[T]
	NewRunTarget(newApi NewApi) Run[T]
}

func (d *Module[T]) NewApi() NewApi         { return di.S(api.New, api.Deps(d.Deps)) }
func (d *Module[T]) WrapApi(api Api) Api    { return api }
func (d *Module[T]) TargetRun() Run[T]      { return d.NewRunTarget(d.NewApi()) }
func (d *Module[T]) RequestTarget() Request { return query.Request.Start }
func (d *Module[T]) FeatOpen() Request      { return open.Feat[T](d) }
