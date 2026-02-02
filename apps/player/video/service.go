package video

import (
	"context"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/ops"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/cryptopunkscc/portal/apps/player"
	"github.com/cryptopunkscc/portal/apps/player/audio"
)

type Service struct {
	audio.Service
	Player player.Video
}

func (s *Service) Serve(ctx context.Context) (err error) {
	if s.Service.Player == nil {
		s.Service.Player = s.Player
	}
	if s.Player == nil {
		s.Player = s.Service.Player.(player.Video)
	}
	if s.ObjectsClient == nil {
		s.ObjectsClient = objects.Default()
	}
	set := ops.NewSet()
	_ = set.AddSubSet("video", ops.Struct(s, "Op"))
	aCtx := astral.NewContext(ctx)
	return ops.Serve(aCtx, set)
}

type opFullScreenArgs struct {
	On int `query:"optional"`
}

func (s *Service) OpFullscreen(_ *astral.Context, query *ops.Query, args opFullScreenArgs) (err error) {
	s.Lock()
	defer s.Unlock()
	conn := query.Accept()
	defer conn.Close()
	on := false
	switch args.On {
	case 1:
		on = true
	case 0:
		if on, err = s.Player.IsFullscreen(); err != nil {
			return
		}
		on = !on
	}
	return s.Player.Fullscreen(on)
}
