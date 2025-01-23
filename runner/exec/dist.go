package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
)

func DistRunner() target.Run[target.DistExec] {
	return func(ctx context.Context, src target.DistExec, args ...string) (err error) {
		abs := src.Target().Executable().Abs()
		return RunCmd(ctx, abs, args...)
	}
}
