package reload

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"strings"
	"time"
)

func Start(
	ctx context.Context,
	portal target.Portal_,
	reload Reload,
	cache api.Cache,
) (send target.MsgSend) {
	c := &client{
		reload:  reload,
		cache:   cache,
		changes: mem.NewCache[time.Time](),
	}
	var err error
	defer plog.PrintTrace(&err)

	if c.conn, err = apphost.Default.Rpc().Conn("portal", "dev.portal.broadcast"); err != nil {
		return
	}
	if err = c.conn.Encode(portal.Manifest().Package); err != nil {
		return
	}

	go c.Handle(ctx)
	send = c.Send
	return
}

type client struct {
	conn rpc.Conn

	reload  Reload
	cache   api.Cache
	changes mem.Cache[time.Time]
}

type Reload func() error

func (s *client) Send(msg target.Msg) (err error) {
	if err = s.conn.Encode(msg); err != nil && err.Error() == "EOF" {
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
	var msg target.Msg
	var err error
	log := plog.Get(ctx).Type(s)
	for {
		if msg, err = rpc.Decode[target.Msg](s.conn); err != nil {
			break
		}
		log.Println("got message", msg)
		s.HandleMsg(ctx, msg)
	}
	if err.Error() != "EOF" {
		plog.Get(ctx).Type(s).F().Println(err)
	}
}

func (s *client) HandleMsg(ctx context.Context, msg target.Msg) {
	log := plog.Get(ctx).D()
	log.Println("received broadcast message:", msg)
	switch msg.Event {
	case target.DevChanged:
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
	case target.DevRefreshed:
		if ok := s.changes.Delete(msg.Pkg); !ok || s.changes.Size() > 0 {
			log.Println("cannot reload", ok, s.changes.Size())
			return
		}
		log.Println("reloading")
		if err := s.reload(); err != nil {
			log.F().Println(err)
		}
	}
}
