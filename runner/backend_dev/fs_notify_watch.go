package backend_dev

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/fsnotify/fsnotify"
	"strings"
)

func fsNotifyWatchWrite(ctx context.Context, path string, file string) (out <-chan any, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	c := make(chan any, 64)
	out = c
	log := plog.Get(ctx).Scope("fsNotifyWatchWrite")
	go func() {
		defer close(c)
		c <- struct{}{}
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) && strings.HasSuffix(event.Name, file) {
					c <- struct{}{}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)

	return
}
