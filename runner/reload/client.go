package reload

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	apphost3 "github.com/cryptopunkscc/portal/factory/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
	apphost2 "github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
)

type client struct {
	flow    rpc.Conn
	handler *handler
}

func newClient() (sender *client) {
	sender = &client{}
	return
}

func (s *client) Init(reRun ReRun, cache apphost.Cache) *client {
	s.handler = newHandler(reRun, cache)
	return s
}

func (s *client) Connect(ctx context.Context, portal target.Portal_) (err error) {
	if s.flow != nil {
		return
	}
	if s.flow, err = apphost2.Rpc(apphost3.Full(ctx)).Client(id.Anyone, "dev.portal.broadcast"); err != nil {
		return
	}
	if err = s.flow.Encode(portal.Manifest().Package); err != nil {
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
		_ = s.flow.Close()
	}()
	var msg target.Msg
	var err error
	log := plog.Get(ctx).Type(s)
	for {
		if msg, err = rpc.Decode[target.Msg](s.flow); err != nil {
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
	if err = s.flow.Encode(msg); err != nil && err.Error() == "EOF" {
		_ = s.Close()
		return
	}
	return
}

func (s *client) Close() (err error) {
	err = s.flow.Close()
	s.flow = nil
	return
}
