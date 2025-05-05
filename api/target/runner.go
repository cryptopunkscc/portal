package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Runner[T any] interface {
	Run(ctx context.Context, src T, args ...string) (err error)
}

type Reloader interface {
	Reload() error
}

var _ ReRunner[Source] = &SourceRunner[Portal_]{}

type SourceRunner[T Portal_] struct {
	Resolve[T]
	Runner[T]
	Reloader
}

func (r *SourceRunner[T]) Try(src Source) (Source, error) {
	return r.Runnable(src)
}

func (r *SourceRunner[T]) Runnable(src Source) (Runnable, error) {
	source, err := r.Resolve(src)
	if err != nil {
		return nil, err
	}
	if any(source) == nil {
		return nil, ErrNotTarget
	}
	return &runnable[T]{
		Portal_:      source,
		source:       source,
		SourceRunner: r,
	}, nil
}

func (r *SourceRunner[T]) setupReloader() {
	if r.Reloader != nil {
		return
	} else if v, ok := r.Runner.(Reloader); ok {
		r.Reloader = v
	} else {
		rr := reRunner[T]{run: r.Runner.Run}
		r.Runner = rr
		r.Reloader = rr
	}
}

func (r *SourceRunner[T]) Run(ctx context.Context, src Source, args ...string) (err error) {
	defer plog.TraceErr(&err)
	r.setupReloader()
	t, err := r.Resolve(src)
	if err != nil {
		return
	}
	return r.Runner.Run(ctx, t, args...)
}

type Runnable interface {
	Portal_
	Source() Portal_
	Run(ctx context.Context, args ...string) error
}

type runnable[T Portal_] struct {
	Portal_
	source T
	*SourceRunner[T]
}

func (r *runnable[T]) Source() Portal_ {
	return r.Portal_
}

var _ Source = &runnable[Portal_]{}
var _ Runnable = &runnable[Portal_]{}

func (r *runnable[T]) Run(ctx context.Context, args ...string) (err error) {
	r.setupReloader()
	return r.Runner.Run(ctx, r.source, args...)
}
