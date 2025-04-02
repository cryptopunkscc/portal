package watcher

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"github.com/fsnotify/fsnotify"
	"time"
)

type reRunner[T target.Dist_] struct {
	reload func(...string) error
	args   []string
}

func ReRunner[T target.Dist_](reload func(...string) error) target.ReRunner[T] {
	return &reRunner[T]{reload: reload}
}

func (r *reRunner[T]) Run(ctx context.Context, dist T, args ...string) (err error) {
	changes, err := fs2.NotifyWatch(ctx, dist.Abs(), fsnotify.Write)
	if err != nil {
		return
	}
	changes = flow.From(changes).Debounce(200 * time.Millisecond)
	r.args = args
	for range changes {
		if err = r.Reload(); err != nil {
			return
		}
	}
	return
}

func (r *reRunner[T]) Reload() error {
	return r.reload(r.args...)
}
