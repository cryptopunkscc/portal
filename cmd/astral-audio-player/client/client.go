package astral_audio_player

import (
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	"github.com/cryptopunkscc/portal/cmd/astral-audio-player/src"
)

type Client struct {
	*astrald.Client
}

func (c Client) PlayID(ctx *astral.Context, objectID astral.ObjectID) (err error) {
	conn, err := c.Query(ctx, "audio.play", query.Args{"id": objectID.String()})
	if err != nil {
		return
	}
	_, _ = conn.Read([]byte{})
	_ = conn.Close()
	return
}

func (c Client) PlayPath(ctx *astral.Context, path string) (err error) {
	conn, err := c.Query(ctx, "audio.play", query.Args{"path": path})
	if err != nil {
		return
	}
	_, _ = conn.Read([]byte{})
	_ = conn.Close()
	return
}

func (c Client) Pause(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, "audio.pause", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Resume(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, "audio.resume", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Stop(ctx *astral.Context) (err error) {
	conn, err := c.Query(ctx, "audio.stop", nil)
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Status(ctx *astral.Context) (status *astral_audio_player.Status, err error) {
	conn, err := c.QueryChannel(ctx, "audio.status", nil)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Switch(channel.Expect(&status), channel.PassErrors)
	err = conn.Send(&astral.Ack{})
	return
}

func (c Client) Seek(ctx *astral.Context, duration time.Duration) (err error) {
	conn, err := c.Query(ctx, "audio.seek", query.Args{"duration": astral.Duration(duration)})
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}

func (c Client) Move(ctx *astral.Context, duration time.Duration) (err error) {
	conn, err := c.Query(ctx, "audio.move", query.Args{"duration": astral.Duration(duration)})
	if err != nil {
		return
	}
	_ = conn.Close()
	return
}
