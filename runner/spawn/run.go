package spawn

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"sync"
)

func NewRunner[T target.Portal](
	wait *sync.WaitGroup,
	find target.Find[T],
	run target.Run[T],
) target.Spawn {
	return func(ctx context.Context, src string) (err error) {
		portals, err := find(src)
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
				log.Println("Runner", t.Abs(), "exit")
			}(t)
		}
		return
	}
}
