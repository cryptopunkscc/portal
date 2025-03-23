package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"os"
)

type Astrald struct {
	exec.Cmd
	NodeRoot mem.String
}

var _ astrald.Runner = &Astrald{}

// Start astral daemon process in a given [context.Context]
func (a *Astrald) Start(ctx context.Context) (err error) {
	cmd := []string{"astrald"}
	if !a.NodeRoot.IsZero() {
		cmd = append(cmd,
			"-root", a.NodeRoot.Require(),
			"-dbroot", a.NodeRoot.Require(),
		)
	}
	a.Cmd.Stdout = os.Stdout
	a.Cmd.Stderr = os.Stderr
	a.Cmd.Env = os.Environ()
	c := a.Cmd.ParseUnsafe(cmd...).Build()
	if err = c.Start(); err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		_ = c.Process.Kill()
	}()
	return
}
