package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"os"
	"time"
)

type Astrald struct {
	NodeRoot string
	exec.Cmd
}

func (a *Astrald) Root() string { return a.NodeRoot }

// Start astral daemon process in a given [context.Context]
func (a *Astrald) Start(ctx context.Context) (err error) {
	cmd := []string{"astrald"}
	if a.NodeRoot != "" {
		cmd = append(cmd,
			"-root", a.NodeRoot,
			"-dbroot", a.NodeRoot,
		)
	}
	a.Cmd.Stdout = os.Stdout
	a.Cmd.Stderr = os.Stderr
	a.Cmd.Env = os.Environ()
	c := a.Cmd.ParseUnsafe(cmd...).Build()
	go func() {
		time.Sleep(25 * time.Millisecond)
		<-ctx.Done()
		_ = c.Process.Kill()
	}()
	return c.Start()
}
