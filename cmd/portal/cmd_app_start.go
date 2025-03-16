package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
)

func (a Application) startApp(ctx context.Context, opt *apphost.PortaldOpenOpt, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting app", cmd)
	return a.Portal.Open(opt, cmd...)
}
