package astral_yt_dlp

import (
	"fmt"
	"io"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/astrald"
)

type Client struct {
	*astrald.Client
}

func (c Client) Download(ctx *astral.Context, request Request) (err error) {
	ch, err := c.QueryChannel(ctx, "yt-dlp.download", request)
	if err != nil {
		return
	}
	return ch.Switch(channel.ExpectAck, channel.PassErrors, channel.WithContext(ctx))
}

func (c Client) Status(ctx *astral.Context) (out <-chan *Progress, e *error) {
	qch, err := c.QueryChannel(ctx, "yt-dlp.status", nil)
	if err != nil {
		return nil, &err
	}
	pch := make(chan *Progress)
	e = new(error)
	go func() {
		defer close(pch)
		err = qch.Switch(
			channel.Chan[*Progress](pch),
			channel.StopOnEOS,
			channel.PassErrors,
			channel.WithContext(ctx),
		)
		if err != nil {
			*e = err
		}
	}()
	return pch, e
}

type Request struct {
	Url   astral.String16
	Dir   astral.String16 `query:"optional"`
	Audio astral.Bool     `query:"optional"` // Download audio instead of video
}

func (r *Request) ObjectType() string { return "yt-dlp.request" }

func (r *Request) String() string { return fmt.Sprint(*r) }

func (r *Request) WriteTo(writer io.Writer) (n int64, err error) {
	return astral.Struct(r).WriteTo(writer)
}

func (r *Request) ReadFrom(reader io.Reader) (n int64, err error) {
	return astral.Struct(r).ReadFrom(reader)
}

type Progress struct {
	*Request
	Percent    astral.Int8
	Speed      astral.String8
	ETA        astral.String8
	Downloaded astral.String8
	Total      astral.String8
	Status     astral.String8
}

func (p *Progress) String() string { return fmt.Sprint(*p) }

func (p *Progress) ObjectType() string { return "yt-dlp.progress" }

func (p *Progress) WriteTo(writer io.Writer) (n int64, err error) {
	return astral.Struct(p).WriteTo(writer)
}

func (p *Progress) ReadFrom(reader io.Reader) (n int64, err error) {
	return astral.Struct(p).ReadFrom(reader)
}

func (p *Progress) UnmarshalText(text []byte) error {
	line := string(text)
	parts := strings.Split(line[5:], "|")
	if len(parts) < 5 {
		return fmt.Errorf("invalid progress line: %s", line)
	}
	var percent float64
	_, _ = fmt.Sscanf(parts[0], "%f", &percent)
	p.Percent = astral.Int8(percent)
	p.Speed = astral.String8(parts[3])
	p.ETA = astral.String8(parts[4])
	p.Downloaded = astral.String8(parts[1])
	p.Total = astral.String8(parts[2])
	p.Status = "downloading"
	return nil
}

func init() {
	if err := astral.Add(
		&Request{},
		&Progress{},
	); err != nil {
		panic(err)
	}
}
