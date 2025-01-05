package srv

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/request"
	"github.com/cryptopunkscc/portal/mock/appstore"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"sync"
)

type Module[T Portal_] struct {
	Deps[T]
	CancelFunc context.CancelFunc
	wg         sync.WaitGroup
	processes  sig.Map[string, T]
	targets    Cache[T]
}

type Deps[T Portal_] interface {
	Run() Run[T]
	Resolve() Resolve[T]
	Priority() Priority
	Handlers() cmd.Handlers
	Processes() *sig.Map[string, T]
}

func (d *Module[T]) Port() apphost.Port                             { return PortPortal }
func (d *Module[T]) Open() Run[string]                              { return request.Create[T](d) }
func (d *Module[T]) Shutdown() context.CancelFunc                   { return d.CancelFunc }
func (d *Module[T]) Observe() func(context.Context, rpc.Conn) error { return appstore.Observe }
func (d *Module[T]) Handlers() cmd.Handlers                         { return cmd.Handlers{} }
func (d *Module[T]) WaitGroup() *sync.WaitGroup                     { return &d.wg }
func (d *Module[T]) Processes() *sig.Map[string, T]                 { return &d.processes }
