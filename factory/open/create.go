package open

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/find"
	"github.com/cryptopunkscc/portal/request/finder"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"sync"
)

type Deps[T target.Portal_] interface {
	WaitGroup() *sync.WaitGroup
	Processes() *sig.Map[string, T]
	Run() target.Run[T]
	Resolve() target.Resolve[T]
	Priority() target.Priority
}

func CreateRun[T target.Portal_](d Deps[T]) target.Run[string] {
	return finder.Requester[T](
		find.Create[T](
			&target.Cache[T]{},
			d.Resolve(),
			d.Priority(),
		),
		supervisor.Runner[T](
			d.WaitGroup(),
			d.Processes(),
			d.Run(),
		),
	)
}
