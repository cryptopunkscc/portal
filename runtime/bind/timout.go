package bind

import (
	"context"
	sig2 "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"time"
)

var ConnectionsThreshold = -1

func ApphostTimeout(ctx context.Context, apphost Apphost, portal target.Portal_) Apphost {
	manifest := portal.Manifest()
	if manifest.Env.Timeout > -1 && ConnectionsThreshold >= 0 {
		go func() {
			duration := 5 * time.Second
			if manifest.Env.Timeout > 0 {
				duration = time.Duration(manifest.Env.Timeout) * time.Millisecond
			}
			t := newTimout(duration, func() {
				_ = sig2.Interrupt()
			})
			t.Enable(true)
			for range apphost.Events().Subscribe(ctx) {
				t.Enable(apphost.Connections().Size() <= ConnectionsThreshold)
			}
		}()
	}
	return apphost
}

type timout struct {
	timeout   time.Duration
	ticker    *time.Ticker
	c         chan any
	onTimeout func()
}

func newTimout(timeout time.Duration, onTimeout func()) *timout {
	return &timout{
		timeout:   timeout,
		onTimeout: onTimeout,
		ticker:    time.NewTicker(timeout),
	}
}

func (t *timout) Enable(value bool) {
	if value {
		go t.Start()
	} else {
		t.Stop()
	}
}

func (t *timout) Start() {
	t.Stop()
	log.Println("timout after", t.timeout)
	t.ticker.Reset(t.timeout)
	select {
	case <-t.c:
	case <-t.ticker.C:
		t.onTimeout()
	}
}

func (t *timout) Stop() {
	if t.c != nil {
		close(t.c)
	}
	t.c = make(chan any)
	t.ticker.Stop()
}
