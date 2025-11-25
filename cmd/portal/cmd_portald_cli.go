package main

import (
	"context"
	"io"
	"strings"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Application) portaldCli(ctx context.Context, cmd ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("running portal cli")

	conn, err := a.Apphost.Query("portald", "cli", nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		_, _ = io.Copy(a, conn)
		cancel()
	}()
	cmd = fixCmd(cmd)
	formatted := strings.Join(cmd, " ") + "\n"
	_, err = conn.Write([]byte(formatted))
	<-ctx.Done()
	return
}
