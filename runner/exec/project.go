package exec

import (
	"context"
	"strings"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	exec2 "github.com/cryptopunkscc/portal/target/exec"
)

func (r Runner) Project() *target.SourceRunner[target.ProjectExec] {
	return &target.SourceRunner[target.ProjectExec]{
		Resolve: target.Any[target.ProjectExec](target.Try(exec2.ResolveProject)),
		Runner:  &ProjectRunner{r},
	}
}

type ProjectRunner struct{ Runner }

func (r *ProjectRunner) Run(ctx context.Context, src target.ProjectExec, args ...string) (err error) {
	cmd := src.Build().Get().Exec
	abs := src.Abs()
	arg := strings.Join(args, " ")
	c, err := exec.Cmd{}.Parse(cmd, abs, arg)
	if err != nil {
		return
	}
	return r.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
}
