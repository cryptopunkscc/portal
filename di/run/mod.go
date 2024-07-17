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
	NewRunTarget(newApi NewApi) Run[T]
	WrapApi(api Api) Api
	TargetFinder() Finder[T]
}

func (d *Module[T]) Path() Path               { return featApps.Path }
func (d *Module[T]) NewApi() NewApi           { return di.Single(api.New, api.Deps(d)) }
func (d *Module[T]) TargetRun() Run[T]        { return d.NewRunTarget(d.NewApi()) }
func (d *Module[T]) TargetFind() Find[T]      { return di.Single(find.New[T], find.Deps[T](d)) }
func (d *Module[T]) TargetCache() *Cache[T]   { return &d.targets }
func (d *Module[T]) TargetDispatch() Dispatch { return query.NewRunner[T](PortOpen).Start }
func (d *Module[T]) FeatOpen() Dispatch       { return di.Single(open.Inject[T], open.Deps[T](d)) }
