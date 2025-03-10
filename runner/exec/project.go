package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runtime/tokens"
	"strings"
)

var ProjectExecRun target.Run[target.ProjectExec] = projectExecRun

func projectExecRun(ctx context.Context, src target.ProjectExec, args ...string) (err error) {
	token, err := tokens.Repository{}.Get(src.Manifest().Package)
	if err != nil {
		return
	}
	e := src.Manifest().Exec
	c := exec.Cmd{}.Parse(e, src.Abs(), strings.Join(args, " "))
	return Cmd{}.Run(ctx, token.Token.String(), c.Cmd, c.Args...)
}
