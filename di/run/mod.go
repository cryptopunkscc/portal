package run

import (
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/runner/query"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
)

type Module[T Base] struct {
	Deps[T]
	di.Cache
	targets Cache[T]
}

type Deps[T Base] interface {
	api.Deps
	open.Deps[T]
	NewRunTarget(newApi NewApi) Run[T]
}

func (d *Module[T]) NewApi() NewApi           { return di.S(api.New, api.Deps(d.Deps)) }
func (d *Module[T]) WrapApi(api Api) Api      { return api }
func (d *Module[T]) TargetRun() Run[T]        { return d.NewRunTarget(d.NewApi()) }
func (d *Module[T]) TargetCache() *Cache[T]   { return &d.targets }
func (d *Module[T]) TargetDispatch() Dispatch { return query.NewOpen().Start }
func (d *Module[T]) FeatOpen() Dispatch       { return open.Inject[T](d) }
