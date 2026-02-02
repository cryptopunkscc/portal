package player

import (
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/portal/apps/player"
)

type Client struct {
	Name string
	*astrald.Client
}

func (c Client) PlayID(ctx *astral.Context, objectID astral.ObjectID) (err error) {
	conn, err := c.Query(ctx, c.Name+".play", query.Args{"id": objectID.String()})
	if err != nil {
		return
	}
	_, _ = conn.Read([]byte{})
	_ = conn.Close()
	return
}

func (c Client) PlayPath(ctx *astral.Context, path string) (err error) {
	conn, err := c.Query(ctx, c.Name+".play", query.Args{"path": path})
	if err != nil {
		return
	}
	_, _ = conn.Read([]byte{})
	_ = conn.Close()
	return
}

func (c Client) Pause(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, c.Name+".pause", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Resume(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, c.Name+".resume", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Stop(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, c.Name+".stop", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Status(ctx *astral.Context) (status *player.Status, err error) {
	conn, err := c.QueryChannel(ctx, c.Name+".status", nil)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Switch(channel.Expect(&status), channel.PassErrors)
	err = conn.Send(&astral.Ack{})
	return
}

func (c Client) Seek(ctx *astral.Context, duration time.Duration) (err error) {
	conn, err := c.Query(ctx, c.Name+".seek", query.Args{"duration": astral.Duration(duration)})
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Move(ctx *astral.Context, duration time.Duration) (err error) {
	conn, err := c.Query(ctx, c.Name+".move", query.Args{"duration": astral.Duration(duration)})
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Fullscreen(ctx *astral.Context, on int) (err error) {
	conn, err := c.Query(ctx, c.Name+".fullscreen", query.Args{"on": on})
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}
