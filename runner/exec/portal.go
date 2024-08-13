package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"os"
	"os/exec"
)

type portal[T target.Portal_] struct {
	command []string
	cmd     *exec.Cmd
	ctx     context.Context
	src     T
}

func Portal[T target.Portal_](command ...string) target.Runner[T] {
	return &portal[T]{command: command}
}

func (p *portal[T]) Reload() (err error) {
	if c := p.cmd; c != nil {
		_ = c.Cancel()
	}
	return p.Run(p.ctx, p.src)
}

func (p *portal[T]) Run(ctx context.Context, src T) (err error) {
	p.ctx = ctx
	p.src = src
	cmd := p.command[0]
	args := append(p.command[1:], src.Abs())
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
		err = fmt.Errorf("exec.Portal '%s': %w", p.command, err)
	}
	return
}
