package reload

import (
	"context"
	"strings"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/query"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/dev"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Start(
	ctx context.Context,
	pkg string,
	reload Reload,
	cache api.Cache,
) (send dev.SendMsg) {
	c := &client{
		reload:  reload,
		cache:   cache,
		changes: mem.NewCache[time.Time](),
	}
	var err error
	defer plog.PrintTrace(&err)

	ah := apphost.Default
	portaldId, err := ah.Resolve("portald")
	if err != nil {
		err = nil
		return
	}
	c.conn, err = ah.WithTarget(portaldId).QueryChannel(
		astral.NewContext(ctx),
		"dev.portal.broadcast",
		query.Args{"pkg": pkg},
	)
	if err != nil {
		return
	}

	go c.Handle(ctx)
	send = c.Send
	return
}

type client struct {
	conn *channel.Channel

	reload  Reload
	cache   api.Cache
	changes mem.Cache[time.Time]
}

type Reload func(ctx context.Context) error

func (s *client) Send(msg *dev.Msg) (err error) {
	if err = s.conn.Send(msg); err != nil && err.Error() == "EOF" {
		_ = s.Close()
		return
	}
	return
}

func (s *client) Close() (err error) {
	err = s.conn.Close()
	s.conn = nil
	return
}

func (s *client) Handle(ctx context.Context) {
	go func() {
		<-ctx.Done()
		_ = s.conn.Close()
	}()
	var msg astral.Object
	var err error
	log := plog.Get(ctx).Type(s)
	for {

		if msg, err = s.conn.Receive(); err != nil {
			break
		}
		log.Println("got message", msg)
		s.HandleMsg(ctx, msg)
	}
	if err.Error() != "EOF" {
		plog.Get(ctx).Type(s).F().Println(err)
	}
}

func (s *client) HandleMsg(ctx context.Context, obj astral.Object) {
	log := plog.Get(ctx).D()
	msg := obj.(*dev.Msg)
	log.Println("received broadcast message:", msg)
	switch msg.Event {
	case dev.Changed:
		if s.cache == nil {
			return
		}
		for _, c := range s.cache.Connections().Copy() {
			if c.In() {
				continue
			}
			query := strings.TrimPrefix(c.Query(), "dev.")
			if strings.HasPrefix(query, msg.Pkg) {
				s.changes.Set(msg.Pkg, msg.Time)
			}
		}
	case dev.Refreshed:
		if ok := s.changes.Delete(msg.Pkg); !ok || s.changes.Size() > 0 {
			log.Println("cannot reload", ok, s.changes.Size())
			return
		}
		log.Println("reloading")
		if err := s.reload(ctx); err != nil {
			log.F().Println(err)
		}
	}
}
