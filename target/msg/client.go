package msg

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Client struct {
	port    target.Port
	flow    rpc.Conn
	handler *Handler
}

func NewClient(port target.Port) (sender *Client) {
	sender = &Client{port: port}
	return
}

func (s *Client) Init(reloader Reloader, apphost target.ApphostCache) *Client {
	s.handler = NewHandler(reloader, apphost)
	return s
}

func (s *Client) Connect(ctx context.Context, portal target.Portal) (err error) {
	if s.flow != nil {
		return
	}
	if s.flow, err = rpc.QueryFlow(id.Anyone, s.port.String()); err != nil {
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

func (s *Client) handle(ctx context.Context) {
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

func (s *Client) Send(msg target.Msg) (err error) {
	if err = s.flow.Encode(msg); err != nil && err.Error() == "EOF" {
		_ = s.Close()
		return
	}
	return
}

func (s *Client) Close() (err error) {
	err = s.flow.Close()
	s.flow = nil
	return
}
