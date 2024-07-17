package srv

import (
	"github.com/cryptopunkscc/astrald/sig"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/tray"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/cryptopunkscc/portal/target/spawn"
	"sync"
)

type Module[T Portal] struct {
	Deps[T]
	di.Cache
	wg        sync.WaitGroup
	processes sig.Map[string, T]
	targets   Cache[T]
}

type Deps[T Portal] interface {
	find.Deps[T]
	spawn.Deps[T]
}

func (d *Module[T]) Port() Port                     { return PortPortal }
func (d *Module[T]) Path() Path                     { return featApps.Path }
func (d *Module[T]) NewTray() NewTray               { return tray.NewRun }
func (d *Module[T]) Processes() *sig.Map[string, T] { return &d.processes }
func (d *Module[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Module[T]) TargetFind() Find[T]            { return di.S(find.New[T], find.Deps[T](d.Deps)) }
func (d *Module[T]) TargetCache() *Cache[T]         { return &d.targets }
func (d *Module[T]) RunSpawn() Dispatch             { return di.S(spawn.Inject[T], spawn.Deps[T](d.Deps)) }
func (d *Module[T]) RunService() serve.Service      { return service.NewRun }
func (d *Module[T]) FeatObserve() serve.Observe     { return featApps.Observe }
