package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
)

type Feat[T target.Portal] struct {
	resolve target.Resolve[T]
	run     target.Run[T]
}

func NewFeat[T target.Portal](resolve target.Resolve[T], run target.Run[T]) target.Spawn {
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
