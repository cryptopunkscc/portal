package dev

import (
	"context"
	"time"

	"github.com/cryptopunkscc/portal/pkg/source/app"
	"github.com/cryptopunkscc/portal/pkg/util/flow"
	"github.com/cryptopunkscc/portal/pkg/util/fs2"
	"github.com/fsnotify/fsnotify"
)

type Reloadable interface {
	Reload(ctx context.Context) error
}

func ReloadOnChange(ctx context.Context, reloadable Reloadable, dist app.Dist) (err error) {
	changes, err := fs2.NotifyWatch(ctx, dist.Path, fsnotify.Write)
	if err != nil {
		return
	}
	for range flow.From(changes).Debounce(200 * time.Millisecond) {
		if err = reloadable.Reload(ctx); err != nil {
			return
		}
	}
	return
}
