package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"time"
)

func Sys(ctx context.Context) (bind.Sys, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &sys{
		log:    plog.Get(ctx).Type(sys{}),
		cancel: cancel,
	}, ctx
}

type sys struct {
	log    plog.Logger
	cancel context.CancelFunc
	exit   int
}

func (s *sys) Log(str string)       { s.log.Scope("Log").Println(str) }
func (s *sys) Sleep(duration int64) { time.Sleep(time.Duration(duration) * time.Millisecond) }
func (s *sys) Exit(code int) {
	s.exit = code
	s.cancel()
}
