package main

import (
	"context"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a Application) startApp(ctx context.Context, opt *portald.OpenOpt, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting app", cmd)
	return a.Portal.Open(opt, cmd...)
}
