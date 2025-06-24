package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/runner/exec"
)

func (s *Service) Connect(ctx context.Context, conn rpc.Conn, opt portald.OpenOpt, args ...string) (err error) {
	ctx = exec.WithReadWriter(ctx, conn)
	if err = s.Open()(ctx, opt, args...); err != nil {
		_ = conn.Encode(err)
	}
	return rpc.Close
}
