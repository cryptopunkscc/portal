package exec

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/target"
	"os"
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
