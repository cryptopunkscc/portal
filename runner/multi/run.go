package multi

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"reflect"
)

type runner[T target.Portal_] struct {
	runners []target.Run[target.Portal_]
}

func Runner[T target.Portal_](
	runners ...target.Run[target.Portal_],
) target.Run[T] {
	return runner[T]{runners: runners}.Run
}

func (r runner[T]) Run(ctx context.Context, portal T, args ...string) (err error) {
	for _, run := range r.runners {
		err = run(ctx, portal, args...)
		if !errors.Is(err, target.ErrNotTarget) {
			return
		}
	}
	return plog.Errorf("invalid target %v: %v", reflect.TypeOf(portal), portal.Path())
}
