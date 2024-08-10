package srv

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/factory/request"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/mock/appstore"
	"github.com/cryptopunkscc/portal/pkg/di"
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
	Run() Run[T]
	Resolve() Resolve[T]
	Priority() Priority
	Handlers() serve.Handlers
	Processes() *sig.Map[string, T]
}

func (d *Module[T]) Port() Port                     { return PortPortal }
func (d *Module[T]) Open() Request                  { return request.Create[T](d) }
func (d *Module[T]) Shutdown() context.CancelFunc   { return d.CancelFunc }
func (d *Module[T]) Observe() serve.Observe         { return appstore.Observe }
func (d *Module[T]) Handlers() serve.Handlers       { return serve.Handlers{} }
func (d *Module[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Module[T]) Processes() *sig.Map[string, T] { return &d.processes }
