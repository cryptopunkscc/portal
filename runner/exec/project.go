package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	exec2 "github.com/cryptopunkscc/portal/resolve/exec"
	"strings"
)

var ProjectRunner = target.SourceRunner[target.ProjectExec]{
	Resolve: target.Any[target.ProjectExec](target.Try(exec2.ResolveProject)),
	Runner:  Project,
}

var Project target.Run[target.ProjectExec] = runProjectExec

func runProjectExec(ctx context.Context, src target.ProjectExec, args ...string) (err error) {
	e := src.Manifest().Exec
	c, err := exec.Cmd{}.Parse(e, src.Abs(), strings.Join(args, " "))
	if err != nil {
		return
	}
	return Cmd{}.RunApp(ctx, *src.Manifest(), c.Cmd, c.Args...)
}
