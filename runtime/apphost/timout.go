package apphost

import (
	"context"
	sig2 "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"time"
)

var ConnectionsThreshold = -1

func WithTimeout(ctx context.Context, apphost target.Apphost, portal target.Portal_) target.Apphost {
	manifest := portal.Manifest()
	if manifest.Env.Timeout > -1 && ConnectionsThreshold >= 0 {
		go func() {
			duration := 5 * time.Second
			if manifest.Env.Timeout > 0 {
				duration = time.Duration(manifest.Env.Timeout) * time.Millisecond
			}
			timeout := NewTimout(duration, func() {
				_ = sig2.Interrupt()
			})
			timeout.Enable(true)
			for range apphost.Events().Subscribe(ctx) {
				activeConnections := len(apphost.Connections()) // TODO optimize
				timeout.Enable(activeConnections <= ConnectionsThreshold)
			}
		}()
	}
	return apphost
}

type Timout struct {
	timeout   time.Duration
	ticker    *time.Ticker
	c         chan any
	onTimeout func()
}

func NewTimout(timeout time.Duration, onTimeout func()) *Timout {
	return &Timout{
		timeout:   timeout,
		onTimeout: onTimeout,
		ticker:    time.NewTicker(timeout),
	}
}

func (t *Timout) Enable(value bool) {
	if value {
		go t.Start()
	} else {
		t.Stop()
	}
}

func (t *Timout) Start() {
	t.Stop()
	log.Println("timout after", t.timeout)
	t.ticker.Reset(t.timeout)
	select {
	case <-t.c:
	case <-t.ticker.C:
		t.onTimeout()
	}
}

func (t *Timout) Stop() {
	if t.c != nil {
		close(t.c)
	}
	t.c = make(chan any)
	t.ticker.Stop()
}
