package contacts

import (
	"context"
	"github.com/cryptopunkscc/astrald/net"
	"github.com/cryptopunkscc/astrald/sig"
	"time"
)

type Link struct {
	Id       int
	RemoteId string
	Remote   string
	Network  string
	Idle     time.Duration
	Since    time.Duration
	Latency  time.Duration
}

func (srv service) Links(ctx context.Context) <-chan []Link {
	c := make(chan []Link)
	go func() {
		c <- srv.links()

		events := srv.node.Network().Events().Subscribe(ctx)
		for range events {
			c <- srv.links()
		}
	}()
	return c
}

func (srv service) links() (contacts []Link) {
	for _, l := range srv.node.Network().Links().All() {
		if l == nil {
			continue
		}
		var idle time.Duration = -1
		var lat time.Duration = -1

		if i, ok := l.Link.(sig.Idler); ok {
			idle = i.Idle().Round(time.Second)
		}

		if l, ok := l.Link.(checkLatency); ok {
			lat = l.Latency()
		}

		c := Link{
			Id:       l.ID(),
			RemoteId: l.RemoteIdentity().String(),
			Remote:   srv.node.Resolver().DisplayName(l.RemoteIdentity()),
			Network:  net.Network(l),
			Idle:     idle,
			Since:    time.Since(l.AddedAt()).Round(time.Second),
			Latency:  lat.Round(time.Millisecond),
		}

		contacts = append(contacts, c)
	}
	return
}

type checkLatency interface {
	Latency() time.Duration
}
