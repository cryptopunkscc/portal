package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/exec"
)

var DistRunner = target.SourceRunner[target.DistExec]{
	Resolve: target.Any[target.DistExec](target.Try(exec.ResolveDist)),
	Runner:  Dist,
}

var Dist target.Run[target.DistExec] = runDist

func runDist(ctx context.Context, src target.DistExec, args ...string) (err error) {
	abs := src.Target().Executable().Abs()
	return Cmd{}.RunApp(ctx, *src.Manifest(), abs, args...)
}
