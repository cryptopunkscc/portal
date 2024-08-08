package finder

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
)

type dispatcher[T target.Portal_] struct {
	find target.Find[T]
	run  target.Run[T]
}

func Dispatcher[T target.Portal_](
	find target.Find[T],
	run target.Run[T],
) target.Dispatch {
	return dispatcher[T]{
		find: find,
		run:  run,
	}.Dispatch
}

func (r dispatcher[T]) Dispatch(ctx context.Context, src string, args ...string) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	portals, err := r.find(ctx, src)
	if err != nil {
		return
	}
	log.D().Printf("found %d portals for %s", len(portals), src)
	for _, t := range portals {
		go func(t T) {
			if err = r.run(ctx, t); err != nil {
				log.E().Println(err)
			}
		}(t)
	}
	return
}
