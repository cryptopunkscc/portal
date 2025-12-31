package portald

import (
	"context"

	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/runner/deprecated/exec"
)

func (s *Service) Connect(ctx context.Context, conn rpc.Conn, opt apphost.OpenOptLegacy, args ...string) (err error) {
	ctx = exec.WithReadWriter(ctx, conn)
	if err = s.Open()(ctx, opt, args...); err != nil {
		_ = conn.Encode(err)
	}
	return rpc.Close
}
