package exec

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/pkg/exec"
)

type Astrald struct {
	exec.Cmd
	NodeRoot string
}

var _ astrald.Runner = &Astrald{}

// Start astral daemon process in a given [context.Context]
func (a *Astrald) Start(ctx context.Context) (err error) {
	cmd := []string{"astrald"}
	if len(a.NodeRoot) > 0 {
		nodeRoot := filepath.ToSlash(a.NodeRoot)
		cmd = append(cmd,
			"-root", nodeRoot,
			"-dbroot", nodeRoot,
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
