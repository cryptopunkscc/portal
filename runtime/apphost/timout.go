package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	sig2 "github.com/cryptopunkscc/portal/pkg/sig"
	"log"
	"time"
)

var ConnectionsThreshold = -1

func Timeout(ctx context.Context, apphost apphost.Cache, portal target.Portal_) {
	if ConnectionsThreshold < 0 {
		return
	}
	timeout := portal.Manifest().Env.Timeout
	if timeout < 0 {
		return
	}
	go func() {
		duration := 5 * time.Second
		if timeout > 0 {
			duration = time.Duration(timeout) * time.Millisecond
		}
		t := newTimout(duration, func() {
			_ = sig2.Interrupt()
		})
		t.Enable(true)
		for e := range apphost.Events().Subscribe(ctx) {
			log.Println("apphost event", e.Type, e.Port, e.Ref)
			t.Enable(apphost.Connections().Size() <= ConnectionsThreshold)
		}
	}()
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
	log.Println("timeout", value)
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
		log.Println("timout!!!")
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
