package astral_audio_player

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/ops"
	"github.com/cryptopunkscc/astrald/mod/media"
	"github.com/gopxl/beep/v2/speaker"
)

func Serve(ctx context.Context) (err error) {
	set := ops.NewSet()
	_ = set.AddSubSet("audio", ops.Struct(Default, "Op"))

	aCtx := astral.NewContext(ctx)
	return ops.Serve(aCtx, set)
}

type OpPlayArgs struct {
	Path string           `query:"optional"`
	ID   *astral.ObjectID `query:"optional"`
}

func (p *Player) OpPlay(ctx *astral.Context, query *ops.Query, args *OpPlayArgs) (err error) {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	if args != nil {
		switch {
		case args.Path != "":
			err = p.SetPath(args.Path)
		case args.ID != nil:
			err = p.SetID(ctx, args.ID)
		}
	}
	return p.Play(ctx)
}

func (p *Player) OpPause(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return speaker.Suspend()
}

func (p *Player) OpResume(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return speaker.Resume()
}

func (p *Player) OpStop(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return p.Close()
}

type opSeekArgs struct {
	Duration string
}

func (p *Player) OpSeek(_ *astral.Context, query *ops.Query, args *opSeekArgs) (err error) {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	duration, err := time.ParseDuration(args.Duration)
	if err != nil {
		return
	}
	return p.Seek(duration)
}

func (p *Player) OpAdd(_ *astral.Context, query *ops.Query, args opSeekArgs) (err error) {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	duration, err := time.ParseDuration(args.Duration)
	if err != nil {
		return
	}
	return p.Move(duration)
}

type opStatusArgs struct {
	Out string `query:"optional"`
}

func (p *Player) OpStatus(_ *astral.Context, query *ops.Query, args opStatusArgs) error {
	p.Lock()
	defer p.Unlock()
	stream := p.steamer
	if stream == nil {
		return query.Reject()
	}

	status := Status{
		Position: astral.Duration(p.CurrentTime()),
		Length:   astral.Duration(p.TotalTime()),
		Track:    p.Track,
	}

	ch := query.AcceptChannel(channel.WithOutputFormat(args.Out))
	defer ch.Close()
	err := ch.Send(&status)
	err = ch.Switch(channel.ExpectAck, channel.PassErrors)
	return err
}

type Status struct {
	Position astral.Duration
	Length   astral.Duration
	Track    *media.AudioFile
}

func (Status) ObjectType() string { return "player.audio.status" }

func (s Status) WriteTo(w io.Writer) (n int64, err error) {
	return astral.Struct(s).WriteTo(w)
}

func (s *Status) ReadFrom(r io.Reader) (n int64, err error) {
	return astral.Struct(s).ReadFrom(r)
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, s)
}

func init() {
	_ = astral.Add(&Status{})
}
