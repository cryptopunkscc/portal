package reload

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

type client struct {
	conn    rpc.Conn
	handler *handler
}

func newClient() (sender *client) {
	sender = &client{}
	return
}

func (s *client) Init(reRun Reload, cache api.Cache) *client {
	s.handler = newHandler(reRun, cache)
	return s
}

func (s *client) Connect(ctx context.Context, portal target.Portal_) (err error) {
	if s.conn != nil {
		return
	}
	if s.conn, err = apphost.Default.Rpc().Client("portal", "dev.portal.broadcast"); err != nil {
		return
	}
	if err = s.conn.Encode(portal.Manifest().Package); err != nil {
		return
	}
	if s.handler != nil {
		go s.handle(ctx)
	}
	return
}

func (s *client) handle(ctx context.Context) {
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
		s.handler.HandleMsg(ctx, msg)
	}
	if err.Error() != "EOF" {
		plog.Get(ctx).Type(s).F().Println(err)
	}
}

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
