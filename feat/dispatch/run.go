package dispatch

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	osexec "os/exec"
	"strings"
	"time"
)

type Feat struct {
	prefix     []string
	executable string
}

func NewFeat(executable string, prefix ...string) target.Dispatch {
	return Feat{executable: executable, prefix: prefix}.Run
}

func (f Feat) Run(
	ctx context.Context,
	src string,
	_ ...string,
) (err error) {

	if err = f.queryOpen(ctx, src); err == nil {
		return
	}

	if err = f.portalServe(ctx); err != nil {
		return
	}

	if err = exec.Retry(8*time.Second, func(i int, n int, duration time.Duration) error {
		return f.queryOpen(ctx, src)
	}); err != nil {
		return
	}

	return
}

func (f Feat) queryOpen(ctx context.Context, src string) (err error) {
	port := strings.Join(append(f.prefix, "portal.open"), ".")
	var conn rpc.Conn
	if conn, err = rpc.QueryFlow(id.Anyone, port); err != nil {
		err = fmt.Errorf("Feat.queryOpen %s: %v", port, err)
		return
	}
	err = rpc.Command(conn, "", src)
	if err != nil {
		return

	}
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	return
}

func (f Feat) portalServe(ctx context.Context) (err error) {
	c := osexec.CommandContext(ctx, f.executable, "s", "-t")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
