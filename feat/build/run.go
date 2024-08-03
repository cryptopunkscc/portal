package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/resolve/exec"
	golang "github.com/cryptopunkscc/portal/resolve/go"
	"github.com/cryptopunkscc/portal/resolve/html"
	"github.com/cryptopunkscc/portal/resolve/js"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"slices"
)

type Feat[T target.Portal_] struct {
	resolve      target.Resolve[T]
	newRunDist   func([]target.NodeModule) target.Run[target.Project_]
	runPack      target.Run[target.Dist_]
	dependencies []target.NodeModule
}

func NewFeat[T target.Portal_](
	resolve target.Resolve[T],
	newRunDist func([]target.NodeModule) target.Run[target.Project_],
	runPack target.Run[target.Dist_],
	dependencies []target.NodeModule,
) *Feat[T] {
	return &Feat[T]{
		resolve:      resolve,
		newRunDist:   newRunDist,
		runPack:      runPack,
		dependencies: dependencies,
	}
}

func (r Feat[T]) Run(ctx context.Context, dir string) (err error) {
	if err = r.Dist(ctx, dir); err != nil {
		return fmt.Errorf("cannot build portal apps: %w", err)
	}
	if err = r.Pack(ctx, dir, "."); err != nil {
		return fmt.Errorf("cannot bundle portal apps: %w", err)
	}
	return
}

func (r Feat[T]) Dist(ctx context.Context, dir ...string) (err error) {
	file, err := source.File(dir...)
	if err != nil {
		return err
	}
	resolve := target.Any[target.Project_](
		target.Skip("node_modules"),
		target.Try(js.ResolveProject),
		target.Try(html.ResolveProject),
		target.Try(golang.ResolveProject),
	)
	projects := target.List(resolve, file)
	for _, m := range projects {
		if err = r.newRunDist(r.dependencies)(ctx, m); err != nil {
			return fmt.Errorf("build.Dist: %w", err)
		}
	}
	return
}

func (r Feat[T]) Pack(ctx context.Context, base, sub string) (err error) {
	file, err := source.File(base, sub)
	if err != nil {
		return err
	}

	resolve := target.Any[target.Dist_](
		target.Skip("node_modules"),
		target.Try(js.ResolveDist),
		target.Try(html.ResolveDist),
		target.Try(exec.ResolveDist),
	)
	distributions := target.Portals[target.Dist_](target.List(resolve, file))
	slices.Reverse(distributions)
	distributions = distributions.Reduced()

	err = errors.New("no targets found")
	for _, dist := range distributions {
		if err = r.runPack(ctx, dist); err != nil {
			return fmt.Errorf("bundle target %v: %v", dist.Path(), err)
		}
	}
	return
}
