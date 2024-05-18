package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"os"
	"os/exec"
)

type Portal[T target.Portal] struct {
	src []string
}

func NewPortal[T target.Portal](src ...string) *Portal[T] {
	return &Portal[T]{src: src}
}

var _ target.Run[target.Portal] = (&Portal[target.Portal]{}).Run

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

func NewRunnerByName[T target.Portal](executable, name string) target.Run[T] {
	log.Println("NewRunnerByName", name)
	return NewPortal[T](executable, "r", name).Run
}

func NewRunner[T target.Portal](executable string) target.Run[T] {
	return func(ctx context.Context, src T) (err error) {
		switch v := any(src).(type) {
		case target.Project:
			switch {
			case v.Type().Is(target.TypeFrontend):
				return NewRunnerByName[target.Project](executable, "wails_dev")(ctx, v)
			case v.Type().Is(target.TypeBackend):
				return NewRunnerByName[target.Project](executable, "goja_dev")(ctx, v)
			}
		case target.App:
			switch {
			case v.Type().Is(target.TypeFrontend):
				return NewRunnerByName[target.App](executable, "wails")(ctx, v)
			case v.Type().Is(target.TypeBackend):
				return NewRunnerByName[target.App](executable, "goja")(ctx, v)
			}
		}
		return
	}
}
