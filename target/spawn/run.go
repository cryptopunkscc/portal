package spawn

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"sync"
)

type Runner[T target.Portal] struct {
	wait      *sync.WaitGroup
	find      target.Find[T]
	run       target.Run[T]
	processes sig.Map[string, any]
}

func NewRunner[T target.Portal](
	wait *sync.WaitGroup,
	find target.Find[T],
	run target.Run[T],
) *Runner[T] {
	return &Runner[T]{
		wait: wait,
		find: find,
		run:  run,
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
			r.runPortal(ctx, t)
		}
	}
	return
}

func (r *Runner[T]) runPortal(ctx context.Context, t T) {
	id := t.Manifest().Package
	if _, ok := r.processes.Set(id, 0); !ok {
		return
	}
	r.wait.Add(1)
	go func(t T) {
		defer r.wait.Done()
		defer r.processes.Delete(id)
		log := plog.Get(ctx)
		if err := r.run(ctx, t); err != nil {
			log.Println(err)
		}
		log.Printf("exit %T %s", t, t.Abs())
	}(t)
	return
}
