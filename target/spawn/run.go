package spawn

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"sync"
)

type Deps[T target.Portal] interface {
	WaitGroup() *sync.WaitGroup
	TargetFind() target.Find[T]
	TargetRun() target.Run[T]
	Processes() *sig.Map[string, T]
}

func Inject[T target.Portal](deps Deps[T]) target.Dispatch {
	return NewRunner[T](
		deps.WaitGroup(),
		deps.TargetFind(),
		deps.TargetRun(),
		deps.Processes(),
	).Run
}

type Runner[T target.Portal] struct {
	wait      *sync.WaitGroup
	find      target.Find[T]
	run       target.Run[T]
	processes *sig.Map[string, T]
}

func NewRunner[T target.Portal](
	wait *sync.WaitGroup,
	find target.Find[T],
	run target.Run[T],
	processes *sig.Map[string, T],
) *Runner[T] {
	return &Runner[T]{
		wait:      wait,
		find:      find,
		run:       run,
		processes: processes,
	}
}

func (r *Runner[T]) Run(ctx context.Context, src string, args ...string) (err error) {
	typ := target.ParseType(target.TypeAny, args...)
	log := plog.Get(ctx).Type(r).Set(&ctx)
	portals, err := r.find(ctx, src)
	if err != nil {
		return
	}
	log.D().Printf("found %d portals for %s", len(portals), src)
	for _, t := range portals {
		if t.Type().Is(typ) {
			r.start(ctx, log, t)
		}
	}
	return
}

func (r *Runner[T]) start(ctx context.Context, log plog.Logger, portal T) {
	id := portal.Manifest().Package
	if _, ok := r.processes.Set(id, portal); !ok {
		return
	}
	r.wait.Add(1)
	go func(t T) {
		log.Printf("start %T %s %s", portal, portal.Manifest().Package, portal.Abs())
		defer log.Printf("exit %T %s %s", portal, portal.Manifest().Package, portal.Abs())
		defer r.wait.Done()
		defer r.processes.Delete(id)
		if err := r.run(ctx, t); err != nil {
			log.Println(err)
		}
	}(portal)
	return
}
