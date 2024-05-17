package runner

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"sync"
)

func NewSpawner[T target.Portal](
	wait *sync.WaitGroup,
	resolve target.Resolve[T],
	run target.Run[T],
) func(context.Context, string) error {
	return func(ctx context.Context, src string) (err error) {
		portals, err := resolve(src)
		if err != nil {
			return
		}
		for _, t := range portals {
			wait.Add(1)
			go func(t T) {
				defer wait.Done()
				if err = run(ctx, t); err != nil {
					log.Println(err)
				}
			}(t)
		}
		return
	}
}
