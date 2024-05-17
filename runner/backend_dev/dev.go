package backend_dev

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"
)

type Event uint

const (
	EventReload = Event(iota + 1)
)

func Dev(ctx context.Context, backend Backend, file string, output func(Event)) (err error) {
	if err = backend.Run(file); err != nil {
		return fmt.Errorf("failed to run %s %v", file, err)
	}
	changes, err := fsNotifyWatchWrite(ctx, file, path.Base(file))
	if err != nil {
		return fmt.Errorf("failed to observe changes %s %v", file, err)
	}
	changes = debounce[any](changes, 200*time.Millisecond)
	for range changes {
		log.Println("backend changed", file)
		if err = backend.Run(file); err != nil {
			log.Printf("failed to rerun %s %v", file, err)
		}
		if output != nil {
			output(EventReload)
		}
	}
	return
}
