package dispatch

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/di/find"
	"github.com/cryptopunkscc/portal/dispatch/finder"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"github.com/cryptopunkscc/portal/target"
	"sync"
)

type Deps[T target.Portal_] interface {
	WaitGroup() *sync.WaitGroup
	Processes() *sig.Map[string, T]
	Run() target.Run[T]
	Resolve() target.Resolve[T]
	Priority() target.Priority
}

func Create[T target.Portal_](d Deps[T]) target.Dispatch {
	return finder.Dispatcher[T](
		find.Create[T](
			&target.Cache[T]{},
			d.Resolve(),
			d.Priority(),
		),
		supervisor.NewRun[T](
			d.WaitGroup(),
			d.Processes(),
			d.Run(),
		),
	)
}
