package apphost

import (
	"log"
	"time"
)

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
	<-t.ticker.C
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
}
