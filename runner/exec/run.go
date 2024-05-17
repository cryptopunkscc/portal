package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/target"
	"log"
	"os"
	"os/exec"
)

func NewRunner[T target.Portal](executable string) target.Run[T] {
	return func(ctx context.Context, src T) (err error) {
		switch any(src).(type) {
		case target.ProjectFrontend:
			return NewRunnerByName[target.Portal](executable, "wails_dev")(ctx, src)
		case target.ProjectBackend:
			return NewRunnerByName[target.Portal](executable, "goja_dev")(ctx, src)
		case target.AppFrontend:
			return NewRunnerByName[target.Portal](executable, "wails")(ctx, src)
		case target.AppBackend:
			return NewRunnerByName[target.Portal](executable, "goja")(ctx, src)
		}
		return
	}
}

func NewRunnerByName[T target.Portal](executable, name string) target.Run[T] {
	log.Println("NewRunnerByName", name)
	return NewPortal[T](executable, "o", name).Run
}

type Portal[T target.Portal] struct {
	src []string
}

var _ target.Run[target.Portal] = (&Portal[target.Portal]{}).Run

func NewPortal[T target.Portal](src ...string) *Portal[T] {
	return &Portal[T]{src: src}
}

func (p *Portal[T]) Run(ctx context.Context, src T) (err error) {
	cmd := p.src[0]
	args := append(p.src[1:], src.Abs())
	var c *exec.Cmd
	if ctx != nil {
		c = exec.CommandContext(ctx, cmd, args...)
	} else {
		c = exec.Command(cmd, args...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		err = fmt.Errorf("exec.Portal '%s': %w", p.src, err)
	}
	return
}
