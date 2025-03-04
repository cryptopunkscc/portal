package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/client/portald"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runtime/rpc"
	"log"
)

func (s *Runner[T]) Connect(ctx context.Context, conn rpc.Conn, opt portald.OpenOpt, args ...string) (err error) {
	ctx = exec.WithReadWriter(ctx, conn)
	log.Printf("conn: %T\n\n", conn)
	if err = s.Open()(ctx, opt, args...); err != nil {
		_ = conn.Encode(err)
	}
	return rpc.Close
}
