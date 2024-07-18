package run

import (
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/runner/query"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
	"github.com/cryptopunkscc/portal/target/find"
)

type Module[T Portal] struct {
	Deps[T]
	di.Cache
	targets Cache[T]
}

type Deps[T Portal] interface {
	api.Deps
	find.Deps[T]
	open.Deps[T]
	NewRunTarget(newApi NewApi) Run[T]
}

func (d *Module[T]) Path() Path               { return featApps.Path }
func (d *Module[T]) NewApi() NewApi           { return di.S(api.New, api.Deps(d.Deps)) }
func (d *Module[T]) WrapApi(api Api) Api      { return api }
func (d *Module[T]) TargetRun() Run[T]        { return d.NewRunTarget(d.NewApi()) }
func (d *Module[T]) TargetFind() Find[T]      { return di.S(find.New[T], find.Deps[T](d.Deps)) }
func (d *Module[T]) TargetCache() *Cache[T]   { return &d.targets }
func (d *Module[T]) TargetDispatch() Dispatch { return query.NewOpen().Start }
func (d *Module[T]) FeatOpen() Dispatch       { return di.S(open.Inject[T], open.Deps[T](d.Deps)) }
