package exec

import (
	"context"
	"os"
	"os/exec"
)

type Dispatch struct {
	executable string
}

func NewDispatch(executable string) *Dispatch {
	return &Dispatch{executable: executable}
}

func (s Dispatch) Start(ctx context.Context, cmd string, args ...string) error {
	args2 := append([]string{cmd}, args...)
	c := exec.CommandContext(ctx, s.executable, args2...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
