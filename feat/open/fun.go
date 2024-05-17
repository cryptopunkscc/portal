package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type Feat[T target.Portal] struct {
	resolve runtime.Resolve[T]
	run     runtime.Run[T]
}

func NewFeat[T target.Portal](resolve runtime.Resolve[T], run runtime.Run[T]) runtime.Spawn {
	return Feat[T]{resolve: resolve, run: run}.Run
}

func (f Feat[T]) Run(ctx context.Context, path string) (err error) {
	portal, err := f.resolve(path)
	if err != nil {
		return errors.New("cannot resolve portal: " + err.Error())
	}
	for _, t := range portal {
		return f.run(ctx, t)
	}
	return errors.New("no target found")
}
