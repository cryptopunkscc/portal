package build

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"log"
	"path/filepath"
)

type Feat struct {
	clean   func(string) error
	runDist target.Run[target.Project_]
	runPack target.Run[target.Dist_]
}

func NewFeat(
	clean func(string) error,
	runDist target.Run[target.Project_],
	runPack target.Run[target.Dist_],
) *Feat {
	return &Feat{
		clean:   clean,
		runDist: runDist,
		runPack: runPack,
	}
}

func (r Feat) Run(ctx context.Context, dir string) (err error) {
	if err = r.clean(dir); err != nil {
		plog.Get(ctx).Type(r).W().Println(err)
	}
	if err = r.Dist(ctx, dir); err == nil {
		log.Println("* build:", dir)
	}
	if err = r.Pack(ctx, dir, "."); err == nil {
		log.Println("* pack:", dir)
	} else {
		plog.Get(ctx).Type(r).W().Printf("build skipped: %v", err)
	}
	return
}

func (r Feat) Dist(ctx context.Context, dir ...string) (err error) {
	if err = run[target.Project_](ctx, r.runDist, dir, target.Match[target.Project_]); err != nil {
		err = fmt.Errorf("build.Dist: %w", err)
	}
	return
}

func (r Feat) Pack(ctx context.Context, dir ...string) (err error) {
	if err = run[target.Dist_](ctx, r.runPack, dir,
		target.Match[target.Dist_],
	); err != nil {
		err = fmt.Errorf("build.Pack: %w", err)
	}
	return
}

func run[T target.Portal_](ctx context.Context, run target.Run[T], dir []string, matchers ...target.Matcher) (err error) {
	projects, err := findIn[T](ctx, dir, matchers...)
	if err != nil {
		return
	}
	if len(projects) == 0 {
		return target.ErrNotFound
	}
	for _, m := range projects {
		if err = run(ctx, m); err != nil {
			return
		}
	}
	return
}

func findIn[T target.Portal_](ctx context.Context, dir []string, matchers ...target.Matcher) ([]T, error) {
	return target.FindByPath(
		source.File,
		sources.Resolver[T]()).
		Reduced(matchers...).
		Call(ctx, filepath.Join(dir...))
}
