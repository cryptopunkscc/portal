package main

import (
	"context"

	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) startApp(ctx context.Context, opt apphost.OpenOpt) (err error) {
	log := plog.Get(ctx)
	log.Println("starting app", opt.App)
	return a.portald().Open(opt)
}
