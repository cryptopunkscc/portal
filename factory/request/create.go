package request

import (
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/find"
	"github.com/cryptopunkscc/portal/request/finder"
	"github.com/cryptopunkscc/portal/runner/supervisor"
	"github.com/cryptopunkscc/portal/runtime/apps"
	"sync"
)

type Deps[T target.Portal_] interface {
	Apps() target.Source
	WaitGroup() *sync.WaitGroup
	Processes() *sig.Map[string, T]
	Run() target.Run[T]
	Resolve() target.Resolve[T]
	Priority() target.Priority
}

func Create[T target.Portal_](d Deps[T]) target.Request {
	return finder.Requester[T](
		find.Create[T](
			apps.Path(d.Apps()),
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
