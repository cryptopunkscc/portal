package reload

import (
	"context"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

func Start(
	ctx context.Context,
	portal target.Portal_,
	reRun Reload,
	cache api.Cache,
) (send target.MsgSend) {
	c := &client{}
	c.handler = newHandler(reRun, cache)
	var err error
	defer func() {
		if err != nil {
			plog.Get(ctx).Println("cannot connect dev.portal.broadcast: %v", err)
		}
	}()
	if c.conn, err = apphost.Default.Rpc().Client("portal", "dev.portal.broadcast"); err != nil {
		return
	}
	if err = c.conn.Encode(portal.Manifest().Package); err != nil {
		return
	}
	if c.handler != nil {
		go c.handle(ctx)
	}
	send = c.Send
	return
}

type client struct {
	conn    rpc.Conn
	handler *handler
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
