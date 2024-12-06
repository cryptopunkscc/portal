package finder

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Requester[T target.Portal_](find target.Find[T], run target.Run[T]) target.Request {
	return func(ctx context.Context, src string) (err error) {
		log := plog.Get(ctx)
		portals, err := find(ctx, src)
		if err != nil {
			return
		}
		log.D().Printf("found %d portals for %s", len(portals), src)
		for _, t := range portals {
			go func(t T) {
				if err = run(ctx, t); err != nil {
					log.E().Println(err)
				}
			}(t)
		}
		return
	}
}
