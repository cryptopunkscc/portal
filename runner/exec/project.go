package exec

import (
	"context"
	"strconv"

	"github.com/cryptopunkscc/portal/api/target"
	exec2 "github.com/cryptopunkscc/portal/target/exec"
	"github.com/google/shlex"
)

func (r Runner) Project() *target.SourceRunner[target.ProjectExec] {
	return &target.SourceRunner[target.ProjectExec]{
		Resolve: target.Any[target.ProjectExec](target.Try(exec2.ResolveProject)),
		Runner:  &ProjectRunner{r},
	}
}

type ProjectRunner struct{ Runner }

func (r *ProjectRunner) Run(ctx context.Context, src target.ProjectExec, args ...string) (err error) {
	cmd, err := shlex.Split(src.Build().Get().Exec)
	if err != nil {
		return
	}
	abs := strconv.Quote(src.Abs())
	args = append(cmd[1:], append([]string{abs}, args...)...)
	return r.RunApp(ctx, *src.Manifest(), cmd[0], args...)
}
