package fs

import (
	"context"
	"github.com/fsnotify/fsnotify"
)

func NotifyWatch(ctx context.Context, path string, filter ...fsnotify.Op) (out <-chan fsnotify.Event, err error) {
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
				if len(filter) == 0 {
					c <- event
					break
				}
				for _, op := range filter {
					if event.Op == op {
						c <- event
						break
					}
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
