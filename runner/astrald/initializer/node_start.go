package initializer

import (
	"context"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (i *Astrald) startAstrald(ctx context.Context) (err error) {
	defer plog.TraceErr(&err)
	i.log.Println("starting astrald...")
	if err = i.Runner.Start(ctx); err != nil {
		return
	}
	i.log.Println("astrald started")
	return
}
