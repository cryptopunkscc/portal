package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"io"
	"os"
	"os/exec"
)

type portal[T target.Portal] struct {
	src []string
}

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

func (r *Runner[T]) newRunBundleExecutable(ctx context.Context, v target.BundleExec) error {
	log := plog.Get(ctx)
	p := v.Executable().Lift().Path()
	log.Println("path", p)
	temp, err := os.CreateTemp(r.cacheDir, p)
	if err != nil {
		return plog.Err(err)
	}
	e := v.Executable()
	file, err := e.Files().Open(e.Path())
	if err != nil {
		return plog.Err(err)
	}
	if err = temp.Chmod(0755); err != nil {
		return plog.Err(err)
	}
	_, err = io.Copy(temp, file)
	if err != nil {
		return plog.Err(err)
	}
	if err = temp.Close(); err != nil {
		return plog.Err(err)
	}
	defer os.Remove(temp.Name())
	_ = file.Close()
	err = newPortal[target.Portal](temp.Name()).run(ctx, v)
	if err != nil {
		return plog.Err(err)
	}
	return nil
}
