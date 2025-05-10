package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"slices"
	"sync"
)

type Dispatcher struct {
	Provider[Runnable]
	Runner[[]Runnable]
}

func (r Dispatcher) Run(ctx context.Context, src string, args ...string) (err error) {
	rr := r.Provide(src)
	return r.Runner.Run(ctx, rr, args...)
}

type CachedRunner[T Portal_] struct {
	Runner[[]Runnable]
	*Cache[T]
}

func (r *CachedRunner[T]) Run(ctx context.Context, runnables []Runnable, args ...string) error {
	for _, rr := range runnables {
		r.Cache.Add(rr.Source().(T))
	}
	return r.Runner.Run(ctx, runnables, args...)
}

var RunSeq Run[[]Runnable] = runSeq

func runSeq(ctx context.Context, runnables []Runnable, args ...string) (err error) {
	for _, r := range runnables {
		if err = r.Run(ctx, args...); err != nil {
			return
		}
	}
	return
}

type AsyncRunner struct {
	*sync.WaitGroup
}

func (r *AsyncRunner) Run(ctx context.Context, runnables []Runnable, args ...string) error {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(runnables))
	if r.WaitGroup != nil {
		r.WaitGroup.Add(len(runnables))
	}
	for _, rr := range runnables {
		go func() {
			a := slices.Clone(args)
			if err := rr.Run(ctx, a...); err != nil {
				plog.Get(ctx).Println(err)
			}
			if r.WaitGroup != nil {
				r.WaitGroup.Done()
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	return nil
}
