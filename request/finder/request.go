package finder

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Requester[T target.Portal_](find target.Find[T], run target.Run[T]) target.Run[string] {
	return func(ctx context.Context, src string, args ...string) (err error) {
		log := plog.Get(ctx)
		log.D().Printf("src: %s, args: %v", src, args)
		portals, err := find(ctx, src)
		if err != nil {
			return
		}
		log.D().Printf("found %d portals for %s", len(portals), src)
		for _, t := range portals {
			go func(t T) {
				if err = run(ctx, t, args...); err != nil {
					log.E().Println(err)
				}
			}(t)
		}
		return
	}
}
