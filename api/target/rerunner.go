package target

import "context"

type ReRunner[T any] interface {
	Runner[T]
	Reloader[T]
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
	return r.Reload()
}

func (r reRunner[T]) Reload() error {
	if r.cancel != nil {
		r.cancel()
	}
	var ctx context.Context
	ctx, r.cancel = context.WithCancel(r.ctx)
	return r.run(ctx, r.src, r.args...)
}
