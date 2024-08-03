package watcher

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/cryptopunkscc/portal/target"
	"github.com/fsnotify/fsnotify"
	"time"
)

type Runner[T target.Dist_] struct {
	reload func() error
}

func NewRunner[T target.Dist_](reload func() error) *Runner[T] {
	return &Runner[T]{reload: reload}
}

func (r *Runner[T]) Run(ctx context.Context, dist T) (err error) {
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

func (r *Runner[T]) Reload() error {
	return r.reload()
}
