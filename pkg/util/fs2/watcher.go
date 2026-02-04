package fs2

import (
	"context"

	"github.com/fsnotify/fsnotify"
)

func NotifyWatch(ctx context.Context, path string, filter fsnotify.Op) (out <-chan fsnotify.Event, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	c := make(chan fsnotify.Event, 64)
	out = c
	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return

			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filter == 0 || filter.Has(event.Op) {
					c <- event
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
			}
		}
	}()

	err = watcher.Add(path)

	return
}
