package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
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

var _ runtime.Run[target.Portal] = (&Portal[target.Portal]{}).Run

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

func NewRunnerByName[T target.Portal](name string) runtime.Run[T] {
	log.Println("NewRunnerByName", name)
	return NewPortal[T]("portal", "r", name).Run
}

func NewRunner[T target.Portal]() runtime.Run[T] {
	return func(ctx context.Context, src T) (err error) {
		switch v := any(src).(type) {
		case target.Project:
			switch {
			case v.Type().Is(target.Frontend):
				return NewRunnerByName[target.Project]("wails_dev")(ctx, v)
			case v.Type().Is(target.Backend):
				return NewRunnerByName[target.Project]("goja_dev")(ctx, v)
			}
		case target.App:
			switch {
			case v.Type().Is(target.Frontend):
				return NewRunnerByName[target.App]("wails")(ctx, v)
			case v.Type().Is(target.Backend):
				return NewRunnerByName[target.App]("goja")(ctx, v)
			}
		}
		return
	}
}
