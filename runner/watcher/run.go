package watcher

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/fsnotify/fsnotify"
	"time"
)

type runner[T target.Dist_] struct {
	reload func() error
}

func Runner[T target.Dist_](reload func() error) target.Runner[T] {
	return &runner[T]{reload: reload}
}

func (r *runner[T]) Run(ctx context.Context, dist T) (err error) {
	changes, err := fs2.NotifyWatch(ctx, dist.Abs(), fsnotify.Write)
	if err != nil {
		return
	}
	changes = flow.From(changes).Debounce(200 * time.Millisecond)

	for range changes {
		if err = r.Reload(); err != nil {
			return
		}
	}
	return
}

func (r *runner[T]) Reload() error {
	return r.reload()
}
