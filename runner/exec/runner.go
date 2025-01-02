package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"os"
	"os/exec"
)

type runner[T target.Portal_] struct {
	command func(T) ([]string, error)
	cmd     *exec.Cmd
	ctx     context.Context
	src     T
}

func Runner[T target.Portal_](command func(T) ([]string, error)) target.Runner[T] {
	return &runner[T]{command: command}
}

func (p *runner[T]) Reload() (err error) {
	if c := p.cmd; c != nil {
		_ = c.Cancel()
	}
	return p.Run(p.ctx, p.src)
}

func (p *runner[T]) Run(ctx context.Context, src T, args ...string) (err error) {
	p.ctx = ctx
	p.src = src
	command, err := p.command(src)
	if err != nil {
		return err
	}
	cmd := command[0]
	args = append(append(command[1:], src.Abs()), args...)
	var c *exec.Cmd
	if ctx != nil {
		c = exec.CommandContext(ctx, cmd, args...)
	} else {
		c = exec.Command(cmd, args...)
	}
	p.cmd = c
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	err = c.Run()
	if err != nil {
		err = fmt.Errorf("exec.Portal '%s': %w", command, err)
	}
	return
}
