package srv

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/apphost"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/request"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/pkg/require"
	"github.com/cryptopunkscc/portal/resolve/source"
	runtime "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/apps"
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
	Handlers() serve.Handlers
	Processes() *sig.Map[string, T]
}

func (d *Module[T]) Apps() (s Source)               { return require.NoErr(source.File(apps.DefaultDir())) }
func (d *Module[T]) Port() Port                     { return PortPortal }
func (d *Module[T]) Open() Request                  { return request.Create[T](d) }
func (d *Module[T]) Client() apphost.Client         { return runtime.Default() }
func (d *Module[T]) Shutdown() context.CancelFunc   { return d.CancelFunc }
func (d *Module[T]) Handlers() serve.Handlers       { return serve.Handlers{} }
func (d *Module[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *Module[T]) Processes() *sig.Map[string, T] { return &d.processes }
