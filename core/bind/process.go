package bind

import (
	"context"
	"time"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func NewProcess(ctx context.Context) (*Process, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Process{
		log:    plog.Get(ctx).Type(Process{}),
		cancel: cancel,
		code:   -1,
	}, ctx
}

type Process struct {
	log    plog.Logger
	cancel context.CancelFunc
	code   int
}

func (s *Process) Code() int            { return s.code }
func (s *Process) Log(str string)       { s.log.Scope("Log").Println(str) }
func (s *Process) Sleep(duration int64) { time.Sleep(time.Duration(duration) * time.Millisecond) }
func (s *Process) Exit(code int) {
	s.code = code
	s.cancel()
}
