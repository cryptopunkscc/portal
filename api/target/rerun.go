package target

import "context"

type ReRunner[T any] interface {
	Run(ctx context.Context, src T, args ...string) (err error)
	ReRun() error
}

func (r Run[T]) ReRunner() ReRunner[T] {
	return reRunner[T]{run: r}
}

type reRunner[T any] struct {
	run    Run[T]
	ctx    context.Context
	cancel context.CancelFunc
	src    T
	args   []string
}

func (r reRunner[T]) Run(ctx context.Context, src T, args ...string) (err error) {
	r.ctx = ctx
	r.src = src
	r.args = args
	return r.ReRun()
}

func (r reRunner[T]) ReRun() error {
	if r.cancel != nil {
		r.cancel()
	}
	var ctx context.Context
	ctx, r.cancel = context.WithCancel(context.Background())
	return r.run(ctx, r.src, r.args...)
}
