package apphost

import (
	"context"
	"errors"
	"io"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Call(
	ctx *astral.Context,
	client astrald.Client,
	method string,
	args any,
) (err error) {
	defer plog.TraceErr(&err)
	conn, err := client.Query(ctx, method, args)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err = conn.Read(make([]byte, 0)); !errors.Is(err, io.EOF) {
		return nil
	}
	return
}

func Receive[T any](
	ctx *astral.Context,
	client astrald.Client,
	method string,
	args any,
) (out T, err error) {
	defer plog.TraceErr(&err)
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
	defer plog.TraceErr(&err)
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
