package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/find"
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
	log := plog.Get(ctx).Type(r).Set(&ctx)
	if err = r.clean(dir); err != nil {
		log.W().Println(err)
	}
	if err = r.Dist(ctx, dir); err != nil {
		log.W().Printf("cannot build portal apps: %w", err)
	}
	if err = r.Pack(ctx, dir, "."); err != nil {
		log.W().Printf("cannot bundle portal apps: %s %w", dir, err)
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
		return errors.New("no targets found")
	}
	for _, m := range projects {
		if err = run(ctx, m); err != nil {
			return
		}
	}
	return
}

func findIn[T target.Portal_](ctx context.Context, dir []string, matchers ...target.Matcher) ([]T, error) {
	return find.ByPath(
		source.File,
		sources.Resolver[T]()).
		Reduced(matchers...).
		Call(ctx, filepath.Join(dir...))
}
