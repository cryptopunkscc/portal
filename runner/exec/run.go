package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
	"os/exec"
)

func NewRun[T target.Portal](executable string, filter ...target.Type) target.Run[T] {
	t := target.TypeNone
	for _, f := range filter {
		t += f
	}
	return func(ctx context.Context, src T) (err error) {
		log := plog.Get(ctx).Scope("exec.Runner").Set(&ctx)
		if t != target.TypeNone && !src.Type().Is(t) {
			log.F().Println(src.Abs(), target.ErrNotTarget)
			return target.ErrNotTarget
		}
		log.Println("target:", src.Abs(), src.Manifest().Package)
		switch any(src).(type) {
		case target.ProjectFrontend:
			return newRunByName[target.Portal](executable, "wails_dev")(ctx, src)
		case target.ProjectBackend:
			return newRunByName[target.Portal](executable, "goja_dev")(ctx, src)
		case target.AppFrontend:
			return newRunByName[target.Portal](executable, "wails")(ctx, src)
		case target.AppBackend:
			return newRunByName[target.Portal](executable, "goja")(ctx, src)
		}
		return
	}
}

func newRunByName[T target.Portal](executable, name string) target.Run[T] {
	return newPortal[T](executable, "o", name).run
}

type portal[T target.Portal] struct {
	src []string
}

var _ target.Run[target.Portal] = (&portal[target.Portal]{}).run

func newPortal[T target.Portal](src ...string) *portal[T] {
	return &portal[T]{src: src}
}

func (p *portal[T]) run(ctx context.Context, src T) (err error) {
	cmd := p.src[0]
	args := append(p.src[1:], src.Abs())
	plog.Get(ctx).Type(src).Printf("%s %v", cmd, args)
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
