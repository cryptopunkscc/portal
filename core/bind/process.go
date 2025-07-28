package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"time"
)

func Process(ctx context.Context) (bind.Process, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &process{
		log:    plog.Get(ctx).Type(process{}),
		cancel: cancel,
		code:   -1,
	}, ctx
}

type process struct {
	log    plog.Logger
	cancel context.CancelFunc
	code   int
}

func (s *process) Code() int            { return s.code }
func (s *process) Log(str string)       { s.log.Scope("Log").Println(str) }
func (s *process) Sleep(duration int64) { time.Sleep(time.Duration(duration) * time.Millisecond) }
func (s *process) Exit(code int) {
	s.code = code
	s.cancel()
}
