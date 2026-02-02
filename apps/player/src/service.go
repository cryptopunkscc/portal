package player

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/astral/channel"
	"github.com/cryptopunkscc/astrald/lib/ops"
	"github.com/cryptopunkscc/astrald/mod/fs"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/cryptopunkscc/portal/apps/player"
)

type Service struct {
	Name string
	sync.Mutex
	player.Player
	ObjectsClient *objects.Client
	Location      *fs.FileLocation
	ObjectID      *astral.ObjectID
}

func (p *Service) Serve(ctx context.Context) (err error) {
	if p.ObjectsClient == nil {
		p.ObjectsClient = objects.Default()
	}

	set := ops.NewSet()
	_ = set.AddSubSet(p.Name, ops.Struct(p, "Op"))

	aCtx := astral.NewContext(ctx)
	return ops.Serve(aCtx, set)
}

func (p *Service) SetID(ctx *astral.Context, id *astral.ObjectID) (err error) {
	p.Location = nil
	p.ObjectID = id
	describeCh, e := p.ObjectsClient.Describe(ctx, id)
	for describe := range describeCh {
		switch d := describe.Descriptor.(type) {
		case *fs.FileLocation:
			p.Location = d
		}
	}
	if e != nil && *e != nil {
		return *e
	}
	return nil
}

func (p *Service) SetPath(filePath string) (err error) {
	p.Location = &fs.FileLocation{Path: astral.String16(filePath)}
	return nil
}

func (p *Service) Play() (err error) {
	var rc io.ReadCloser
	var ext string

	switch {
	case p.Location != nil:
		ext = filepath.Ext(p.Location.Path.String())
		rc, err = os.Open(p.Location.Path.String())
	default:
		return errors.New("no audio file specified")
	}

	return p.Player.Play(rc, ext)
}

type OpPlayArgs struct {
	Path string           `query:"optional"`
	ID   *astral.ObjectID `query:"optional"`
}

func (p *Service) OpPlay(ctx *astral.Context, query *ops.Query, args *OpPlayArgs) (err error) {
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
	return p.Play()
}

func (p *Service) OpPause(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return p.Player.Suspend()
}

func (p *Service) OpResume(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return p.Player.Resume()
}

func (p *Service) OpStop(_ *astral.Context, query *ops.Query) error {
	p.Lock()
	defer p.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return p.Close()
}

type opSeekArgs struct {
	Duration string
}

func (p *Service) OpSeek(_ *astral.Context, query *ops.Query, args *opSeekArgs) (err error) {
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

func (p *Service) OpMove(_ *astral.Context, query *ops.Query, args opSeekArgs) (err error) {
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

func (p *Service) OpStatus(_ *astral.Context, query *ops.Query, args opStatusArgs) error {
	p.Lock()
	defer p.Unlock()

	status := player.Status{
		ObjectID: p.ObjectID,
		Position: astral.Duration(p.CurrentTime()),
		Length:   astral.Duration(p.TotalTime()),
	}

	ch := query.AcceptChannel(channel.WithOutputFormat(args.Out))
	defer ch.Close()
	err := ch.Send(&status)
	err = ch.Switch(channel.ExpectAck, channel.PassErrors)
	return err
}
