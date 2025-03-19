package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	sig2 "github.com/cryptopunkscc/portal/pkg/sig"
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
		log := plog.Get(ctx).D().Type(t).Set(&ctx)
		t.log = log
		t.Enable(true)
		for e := range apphost.Events().Subscribe(ctx) {
			log.Printf("apphost event %v %s %s", e.Type, e.Query, e.Ref)
			t.Enable(apphost.Connections().Size() <= ConnectionsThreshold)
		}
	}()
}

type timout struct {
	log       plog.Logger
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
	t.log.Println("timeout", value)
	t.stop()
	if value {
		go t.start()
	}
}

func (t *timout) start() {
	t.log.Println("timout after", t.timeout)
	t.ticker.Reset(t.timeout)
	t.c = make(chan any)
	select {
	case <-t.c:
	case <-t.ticker.C:
		t.log.Println("timout!!!")
		t.onTimeout()
	}
}

func (t *timout) stop() {
	if t.c != nil {
		close(t.c)
		t.c = nil
		t.ticker.Stop()
	}
}
