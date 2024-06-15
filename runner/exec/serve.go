package exec

import (
	"context"
	"os"
	"os/exec"
)

type Dispatcher struct {
	executable string
}

func NewDispatcher(executable string) *Dispatcher {
	return &Dispatcher{executable: executable}
}

func (s Dispatcher) Dispatch(ctx context.Context, cmd string, args ...string) error {
	args2 := append([]string{cmd}, args...)
	c := exec.CommandContext(ctx, s.executable, args2...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Start()
}
