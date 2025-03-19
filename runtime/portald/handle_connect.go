package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

func (s *Runner[T]) Connect(ctx context.Context, conn rpc.Conn, opt apphost.PortaldOpenOpt, args ...string) (err error) {
	ctx = exec.WithReadWriter(ctx, conn)
	if err = s.Open()(ctx, opt, args...); err != nil {
		_ = conn.Encode(err)
	}
	return rpc.Close
}
