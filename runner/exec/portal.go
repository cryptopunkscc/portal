package exec

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"io"
	"os"
	"os/exec"
	"path"
)

func NewRun[T target.Portal](cacheDir, executable string) target.Run[T] {
	return NewRunner[T](cacheDir, executable).Run
}

type Runner[T target.Portal] struct {
	cacheDir   string
	executable string
}

func NewRunner[T target.Portal](cacheDir string, executable string) *Runner[T] {
	cacheDir = path.Join(cacheDir, "apps", "tmp")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}
	return &Runner[T]{cacheDir: cacheDir, executable: executable}
}

func (r *Runner[T]) Run(ctx context.Context, src T) (err error) {
	log := plog.Get(ctx).Type(r).Set(&ctx)
	log.Printf("start %T %s %s", src, src.Manifest().Package, src.Abs())
	defer log.Printf("exit %T %s %s", src, src.Manifest().Package, src.Abs())
	switch v := any(src).(type) {
	case target.ProjectHtml:
		return newPortal[target.Portal](r.executable, "o", "wails_dev").run(ctx, src)
	case target.ProjectJs:
		return newPortal[target.Portal](r.executable, "o", "goja_dev").run(ctx, src)
	case target.ProjectGo:
		return newPortal[target.Portal](r.executable, "o", "go_dev").run(ctx, src)
	case target.AppHtml:
		return newPortal[target.Portal](r.executable, "o", "wails").run(ctx, src)
	case target.AppJs:
		return newPortal[target.Portal](r.executable, "o", "goja").run(ctx, src)
	case target.DistExec:
		return newPortal[target.Portal](v.Executable().Abs()).run(ctx, v)
	case target.BundleExec:
		return r.newRunBundleExecutable(ctx, v)
	}
	return
}

type portal[T target.Portal] struct {
	src []string
}

func newPortal[T target.Portal](src ...string) *portal[T] {
	return &portal[T]{src: src}
}

func (p *portal[T]) run(ctx context.Context, src T) (err error) {
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
	c.Stdin = os.Stdin
	err = c.Run()
	if err != nil {
		err = fmt.Errorf("exec.Portal '%s': %w", p.src, err)
	}
	return
}

func (r *Runner[T]) newRunBundleExecutable(ctx context.Context, v target.BundleExec) error {
	p := v.Executable().Lift().Path()
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
		return err
	}
	return nil
}
