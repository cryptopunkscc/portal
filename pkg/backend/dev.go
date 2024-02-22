package backend

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
	"time"
)

func Dev(backend Backend, path string, output chan<- Event) (err error) {
	if err = backend.Run(path); err != nil {
		return fmt.Errorf("failed to run %s %v", path, err)
	}
	changes, err := observeChanges(path, "service.js")
	if err != nil {
		return fmt.Errorf("failed to observe changes %s %v", path, err)
	}
	go func() {
		changes = debounce[any](changes, 200*time.Millisecond)
		for range changes {
			if err = backend.Run(path); err != nil {
				log.Printf("failed to rerun %s %v", path, err)
			}
			if output != nil {
				output <- EventReload
			}
		}
	}()
	return
}

type Event uint

const (
	EventReload = Event(iota + 1)
)

func observeChanges(path string, file string) (out <-chan any, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	c := make(chan any, 64)
	out = c
	go func() {
		defer close(c)
		for {
			select {
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
				//c <- err
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)

	return
}

func debounce[T any](in <-chan T, t time.Duration) (out <-chan T) {
	buff := make(chan T, 64)
	go func() {
		for v := range in {
			buff <- v
		}
	}()

	o := make(chan T)
	out = o
	go func() {
		var last *T = nil
		for {
			select {
			case l := <-buff:
				last = &l
			default:
				time.Sleep(t / 2)
				if last == nil {
					continue
				}
				time.Sleep(t / 2)
				o <- *last
				last = nil
			}
		}
	}()
	return
}
