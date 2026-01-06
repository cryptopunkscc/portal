package main

import (
	"context"

	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) startApp(ctx context.Context, opt *apphost.OpenOptLegacy, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("starting app", cmd)
	return a.portald().Open(opt, cmd...)
}
