package main

import (
	"context"
	"io"
	"os"

	"github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) runApp(ctx context.Context, opt *portald.OpenOpt, cmd []string) (err error) {
	log := plog.Get(ctx)
	log.Println("running app", opt, cmd)

	conn, err := a.portald().Connect(opt, cmd...)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		_, _ = io.Copy(a, conn)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		cancel()
	}()
	<-ctx.Done()
	return
}
