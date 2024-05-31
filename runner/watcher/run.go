package watcher

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/flow"
	"github.com/cryptopunkscc/go-astral-js/pkg/fs"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/fsnotify/fsnotify"
	"time"
)

type Runner[T target.Dist] struct {
	reload func() error
}

func NewRunner[T target.Dist](reload func() error) *Runner[T] {
	return &Runner[T]{reload: reload}
}

func (r *Runner[T]) Run(ctx context.Context, dist T) (err error) {
	changes, err := fs.NotifyWatch(ctx, dist.Abs(), fsnotify.Write)
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
