package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
)

var DistRun target.Run[target.DistExec] = dist

func dist(ctx context.Context, src target.DistExec, args ...string) (err error) {
	abs := src.Target().Executable().Abs()
	return Cmd{}.RunApp(ctx, *src.Manifest(), abs, args...)
}
