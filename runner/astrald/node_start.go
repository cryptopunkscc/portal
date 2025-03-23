package astrald

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (r *Runner) startAstrald(ctx context.Context) (err error) {
	defer plog.TraceErr(&err)
	r.log.Println("starting astrald...")
	if err = r.Runner.Start(ctx); err != nil {
		return
	}
	r.log.Println("astrald started")
	return
}
