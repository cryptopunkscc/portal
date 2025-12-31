package exec

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/astrald"
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

		// FIXME: usage of strconv.Quote works fine for now, but whole command should be built more safely.
		cmd = append(cmd,
			"-root", strconv.Quote(nodeRoot),
			"-dbroot", strconv.Quote(nodeRoot),
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
