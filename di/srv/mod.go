package srv

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/find"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/service"
	"github.com/cryptopunkscc/portal/runner/spawn"
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
func (d *Module[T]) Path() Path                     { return featApps.Path }
func (d *Module[T]) TargetFind() Find[T]            { return find.Inject[T](d) }
func (d *Module[T]) TargetFile() File               { return source.File }
func (d *Module[T]) TargetCache() *Cache[T]         { return &d.targets }
func (d *Module[T]) Embed() []Source                { return []Source{} }

// func (d *Module[T]) Embed() []Source                { return []Source{source.Embed(apps.LauncherSvelteFS)} }
