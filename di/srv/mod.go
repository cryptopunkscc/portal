package srv

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/dispatch/service"
	"github.com/cryptopunkscc/portal/dispatch/spawn"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/resolve/source"
	. "github.com/cryptopunkscc/portal/target"
	"sync"
)

type Module[T Portal_] struct {
	Deps[T]
	di.Cache
	CancelFunc context.CancelFunc
	wg         sync.WaitGroup
	processes  sig.Map[string, T]
	targets    Cache[T]
}

type Deps[T Portal_] interface {
	TargetResolve() Resolve[T]
	Priority() Priority
	TargetRun() Run[T]
}

func (d *Module[T]) Port() Port                     { return PortPortal }
func (d *Module[T]) Close() context.CancelFunc      { return d.CancelFunc }
func (d *Module[T]) RunSpawn() Dispatch             { return spawn.Inject[T](d) }
func (d *Module[T]) RunService() serve.Service      { return service.NewRun }
func (d *Module[T]) FeatObserve() serve.Observe     { return featApps.Observe }
func (d *Module[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Module[T]) Processes() *sig.Map[string, T] { return &d.processes }
func (d *Module[T]) TargetFind() Find[T] {
	return FindByPath(source.File, d.TargetResolve()).
		ById(featApps.Path).
		Cached(&d.targets).
		Reduced(d.Priority()...)
}
