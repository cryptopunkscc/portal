package audio

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
	sync.Mutex
	player.Player
	ObjectsClient *objects.Client
	Location      *fs.FileLocation
	ObjectID      *astral.ObjectID
	Name          string
}

func (s *Service) Serve(ctx context.Context) (err error) {
	if s.Name == "" {
		s.Name = "audio"
	}
	if s.ObjectsClient == nil {
		s.ObjectsClient = objects.Default()
	}

	set := ops.NewSet()
	_ = set.AddSubSet(s.Name, ops.Struct(s, "Op"))

	aCtx := astral.NewContext(ctx)
	return ops.Serve(aCtx, set)
}

func (s *Service) SetID(ctx *astral.Context, id *astral.ObjectID) (err error) {
	s.Location = nil
	s.ObjectID = id
	describeCh, e := s.ObjectsClient.Describe(ctx, id)
	for describe := range describeCh {
		switch d := describe.Descriptor.(type) {
		case *fs.FileLocation:
			s.Location = d
		}
	}
	if e != nil && *e != nil {
		return *e
	}
	return nil
}

func (s *Service) SetPath(filePath string) (err error) {
	s.Location = &fs.FileLocation{Path: astral.String16(filePath)}
	return nil
}

func (s *Service) Play() (err error) {
	var rc io.ReadCloser
	var ext string

	switch {
	case s.Location != nil:
		ext = filepath.Ext(s.Location.Path.String())
		rc, err = os.Open(s.Location.Path.String())
	default:
		return errors.New("no audio file specified")
	}

	return s.Player.Play(rc, ext)
}

type OpPlayArgs struct {
	Path string           `query:"optional"`
	ID   *astral.ObjectID `query:"optional"`
}

func (s *Service) OpPlay(ctx *astral.Context, query *ops.Query, args *OpPlayArgs) (err error) {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	if args != nil {
		switch {
		case args.Path != "":
			err = s.SetPath(args.Path)
		case args.ID != nil:
			err = s.SetID(ctx, args.ID)
		}
	}
	return s.Play()
}

func (s *Service) OpPause(_ *astral.Context, query *ops.Query) error {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return s.Player.Suspend()
}

func (s *Service) OpResume(_ *astral.Context, query *ops.Query) error {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return s.Player.Resume()
}

func (s *Service) OpStop(_ *astral.Context, query *ops.Query) error {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	return s.Close()
}

type opSeekArgs struct {
	Duration string
}

func (s *Service) OpSeek(_ *astral.Context, query *ops.Query, args *opSeekArgs) (err error) {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	duration, err := time.ParseDuration(args.Duration)
	if err != nil {
		return
	}
	return s.Seek(duration)
}

func (s *Service) OpMove(_ *astral.Context, query *ops.Query, args opSeekArgs) (err error) {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	duration, err := time.ParseDuration(args.Duration)
	if err != nil {
		return
	}
	return s.Move(duration)
}

type opStatusArgs struct {
	Out string `query:"optional"`
}

func (s *Service) OpStatus(_ *astral.Context, query *ops.Query, args opStatusArgs) error {
	s.Lock()
	defer s.Unlock()

	status := player.Status{
		ObjectID: s.ObjectID,
		Position: astral.Duration(s.CurrentTime()),
		Length:   astral.Duration(s.TotalTime()),
	}

	ch := query.AcceptChannel(channel.WithOutputFormat(args.Out))
	defer ch.Close()
	err := ch.Send(&status)
	err = ch.Switch(channel.ExpectAck, channel.PassErrors)
	return err
}
