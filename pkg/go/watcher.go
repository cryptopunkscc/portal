package golang

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/fsnotify/fsnotify"
	"strings"
)

type Watcher struct {
	filterExt   []string
	abs         string
	projectRoot string
	moduleName  string
	cache       *WatchCache
	watcher     *fsnotify.Watcher
	events      chan fsnotify.Event
}

func NewWatcher() *Watcher {
	return &Watcher{
		filterExt: []string{".go", "go.mod"},
	}
}

func (w *Watcher) Run(ctx context.Context, abs string) (c <-chan fsnotify.Event, err error) {
	log := plog.Get(ctx).Type(w).Set(&ctx)
	log.Println("starting watcher", abs)
	if w.projectRoot, err = findProjectRoot(abs); err != nil {
		return
	}
	if w.moduleName, err = getModuleRoot(w.projectRoot); err != nil {
		return
	}
	if w.watcher, err = fsnotify.NewWatcher(); err != nil {
		return
	}
	if err = w.watcher.Add(abs); err != nil {
		return
	}
	log.Println("starting watcher", abs, w.projectRoot, w.moduleName)
	w.abs = abs
	w.cache = NewWatchCache(w.projectRoot, w.moduleName)
	for s := range w.cache.AddDir(abs) {
		_ = w.watcher.Add(s)
	}
	w.events = make(chan fsnotify.Event)
	c = w.events

	go w.run(ctx)
	return
}

func (w *Watcher) run(ctx context.Context) {
	log := plog.Get(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-w.watcher.Errors:
			if err != nil {
				log.E().Println(err)
			}
			if !ok {
				return
			}
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if !w.validSuffix(event.Name) {
				continue
			}

			add := make(map[string]any)
			remove := add
			switch event.Op {
			case fsnotify.Create:
				_ = w.watcher.Add(event.Name)
				add = w.cache.AddFile(event.Name)
			case fsnotify.Remove:
				_ = w.watcher.Remove(event.Name)
				remove = w.cache.RemoveFile(event.Name)
			case fsnotify.Write:
				remove, add = w.cache.UpdateFile(event.Name)
			case fsnotify.Rename:
				remove, add = w.cache.UpdateFile(event.Name)
			case fsnotify.Chmod:
				remove, add = w.cache.UpdateFile(event.Name)
			}
			for s := range add {
				_ = w.watcher.Add(s)
			}
			for s := range remove {
				_ = w.watcher.Remove(s)
			}

			w.events <- event
		}
	}
}

func (w *Watcher) validSuffix(name string) bool {
	for _, s := range w.filterExt {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}
