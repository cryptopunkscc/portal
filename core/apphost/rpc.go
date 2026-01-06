package apphost

import (
	"context"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	apphost "github.com/cryptopunkscc/portal/core/apphost/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func (a *Adapter) Rpc() rpc.Rpc {
	return &apphost.Rpc{
		Apphost: a,
		Log:     a.Log,
	}
}

var _ rpc.Rpc = &Adapter{}

func (a *Adapter) Format(name string) rpc.Rpc {
	return a.Rpc().Format(name)
}

func (a *Adapter) Conn(target, query string) (rpc.Conn, error) {
	return a.Rpc().Conn(target, query)
}

func (a *Adapter) Request(target string, query ...string) rpc.Conn {
	return a.Rpc().Request(target, query...)
}

func (a *Adapter) Router(handler cmd.Handler) rpc.Router {
	return a.Rpc().Router(handler)
}

func Request[T any](ctx *astral.Context, client astrald.Client, method string, args any) (out T, err error) {
	conn, err := client.QueryChannel(ctx, method, args)
	if err != nil {
		return
	}
	receive, err := conn.Receive()
	if err != nil {
		return
	}
	return receive.(T), nil
}

func List[T any](
	ctx *astral.Context,
	client astrald.Client,
	method string,
	args any,
	cfg ...channel.ConfigFunc,
) (out []T, err error) {
	c, err := GoChan[T](ctx, client, method, args, cfg...)
	if err != nil {
		return
	}
	for t := range c {
		out = append(out, t)
	}
	return
}

func GoChan[T any](
	ctx *astral.Context,
	client astrald.Client,
	method string,
	args any,
	cfg ...channel.ConfigFunc,
) (out <-chan T, err error) {
	if ctx == nil {
		ctx = astral.NewContext(context.Background())
	}
	if ctx.Identity() == nil {
		ctx = ctx.WithIdentity(client.GuestID())
	}
	c, err := client.QueryChannel(ctx, method, args, cfg...)
	if err != nil {
		return
	}
	res := make(chan T)
	go func() {
		defer close(res)
		_ = c.Handle(ctx, func(o astral.Object) {
			switch o := o.(type) {
			case T:
				res <- o
			case *astral.EOS:
				c.Close()
			default:
				c.Close()
			}
		})
	}()
	return res, nil
}
