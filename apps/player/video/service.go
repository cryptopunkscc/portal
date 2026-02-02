package video

import (
	"context"

	"github.com/cryptopunkscc/portal/apps/player/audio"
)

type Service struct {
	audio.Service
}

func (s *Service) Serve(ctx context.Context) (err error) {
	s.Name = "video"
	return s.Service.Serve(ctx)
}
