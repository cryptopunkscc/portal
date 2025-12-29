package exec

import (
	"context"
	"strconv"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/exec"
)

func (r Runner) Dist() *target.SourceRunner[target.DistExec] {
	return &target.SourceRunner[target.DistExec]{
		Resolve: target.Any[target.DistExec](target.Try(exec.ResolveDist)),
		Runner:  &DistRunner{r},
	}
}

type DistRunner struct{ Runner }

func (r *DistRunner) Run(ctx context.Context, src target.DistExec, args ...string) (err error) {
	abs := src.Runtime().Executable().Abs()
	abs = strconv.Quote(abs)
	return r.RunApp(ctx, *src.Manifest(), abs, args...)
}
