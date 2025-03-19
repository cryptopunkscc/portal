package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"strings"
)

var ProjectExecRun target.Run[target.ProjectExec] = projectExecRun

func projectExecRun(ctx context.Context, src target.ProjectExec, args ...string) (err error) {
	t, err := token.Repository{}.Get(src.Manifest().Package)
	if err != nil {
		return
	}
	e := src.Manifest().Exec
	c, err := exec.Cmd{}.Parse(e, src.Abs(), strings.Join(args, " "))
	if err != nil {
		return
	}
	return Cmd{}.Run(ctx, t.Token.String(), c.Cmd, c.Args...)
}
