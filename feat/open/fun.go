package open

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Feat[T target.Portal] struct {
	find target.Find[T]
	run  target.Run[T]
}

func NewFeat[T target.Portal](find target.Find[T], run target.Run[T]) target.Dispatch {
	return Feat[T]{find: find, run: run}.Run
}

func (f Feat[T]) Run(ctx context.Context, path string, _ ...string) (err error) {
	portal, err := f.find(path)
	if err != nil {
		return errors.New("cannot resolve portal: " + err.Error())
	}
	for _, t := range portal {
		return f.run(ctx, t)
	}
	return errors.New("no target found")
}
