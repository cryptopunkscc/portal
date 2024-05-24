package spawn

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"sync"
)

type Runner[T target.Portal] struct {
	wait *sync.WaitGroup
	find target.Find[T]
	run  target.Run[T]
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

func (r Runner[T]) Run(ctx context.Context, src string, _ ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	portals, err := r.find(ctx, src)
	log.D().Printf("found %d portals for %s", len(portals), src)
	if err != nil {
		return
	}
	for _, t := range portals {
		r.wait.Add(1)
		go func(t T) {
			defer r.wait.Done()
			if err = r.run(ctx, t); err != nil {
				log.Println(err)
			}
			log.Println("Runner", t.Abs(), "exit")
		}(t)
	}
	return
}
