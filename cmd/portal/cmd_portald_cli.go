package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"io"
	"os"
	"strings"
)

func (a Application) portaldCli(ctx context.Context, cmd ...string) (err error) {
	log := plog.Get(ctx)
	log.Println("running portal cli")

	conn, err := apphost.Default.Query("portal", "cli", nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		cancel()
	}()
	cmd = fixCmd(cmd)
	formatted := strings.Join(cmd, " ") + "\n"
	_, err = conn.Write([]byte(formatted))
	<-ctx.Done()
	return
}
