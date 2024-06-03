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
	log := plog.Get(ctx).Scope("exec.Runner").Set(&ctx)
	log.Printf("target: %T %s %s", src, src.Abs(), src.Manifest().Package)
	switch v := any(src).(type) {
	case target.ProjectFrontend:
		return newPortal[target.Portal](r.executable, "o", "wails_dev").run(ctx, src)
	case target.ProjectBackend:
		return newPortal[target.Portal](r.executable, "o", "goja_dev").run(ctx, src)
	case target.AppFrontend:
		return newPortal[target.Portal](r.executable, "o", "wails").run(ctx, src)
	case target.AppBackend:
		return newPortal[target.Portal](r.executable, "o", "goja").run(ctx, src)
	case target.DistExecutable:
		return newPortal[target.Portal](v.Exec().Abs()).run(ctx, v)
	case target.BundleExecutable:
		return r.newRunBundleExecutable(ctx, v)
	}
	return
}
