package backend_dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"path"
	"time"
)

type Event uint

const (
	EventReload = Event(iota + 1)
)

func Dev(ctx context.Context, backend Backend, dist target.Dist, output func(Event)) (err error) {
	if !path.IsAbs(dist.Abs()) {
		return plog.Errorf("cannot run dist with non-absolute path: %s", dist.Abs())
	}
	file := path.Join(dist.Abs(), "main.js")

	log := plog.Get(ctx)
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
