package bind

import (
	"context"
	"github.com/cryptopunkscc/portal/api/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"time"
)

func Sys(ctx context.Context) bind.Sys { return &sys{log: plog.Get(ctx)} }

type sys struct{ log plog.Logger }

func (s *sys) Log(str string)       { s.log.Println(str) }
func (s *sys) Sleep(duration int64) { time.Sleep(time.Duration(duration) * time.Millisecond) }
