package astral_yt_dlp

import (
	"log"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/ops"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/api"
	"github.com/kkdai/youtube/v2"
)

type Service struct {
	youtube.Client
	queue    Queue
	download Download
	Dir      string
}

func (s *Service) Serve(ctx *astral.Context) (err error) {
	set := ops.NewSet()
	_ = set.AddSubSet("yt-dlp", ops.Struct(s, "Op"))
	go func() {
		for request := range sig.Subscribe(ctx, &s.queue.Queue) {
			if err = s.download.Run(ctx, request); err != nil {
				log.Printf("Error downloading %v: %v", request, err)
			}
			if s.queue.Done(request) {
				s.download.Reset()
			}
		}
	}()
	return ops.Serve(ctx, set)
}

func (s *Service) OpDownload(_ *astral.Context, query *ops.Query, request astral_yt_dlp.Request) (err error) {
	if len(request.Dir) == 0 {
		request.Dir = astral.String16(s.Dir)
	}
	ch := query.AcceptChannel()
	defer ch.Close()
	if err = s.queue.Push(request); err != nil {
		return ch.Send(astral.Err(err))
	}
	return ch.Send(&astral.Ack{})
}

func (s *Service) OpStatus(ctx *astral.Context, query *ops.Query) (err error) {
	ch := query.AcceptChannel()
	defer ch.Close()
	for progress := range s.download.Progress(ctx.Context, false) {
		err = ch.Send(&progress)
		if err != nil {
			return err
		}
	}
	return
}
