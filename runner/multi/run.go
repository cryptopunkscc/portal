package multi

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"reflect"
)

type Runner[T target.Portal] struct {
	runners []target.Run[target.Portal]
}

func NewRunner[T target.Portal](
	runners ...target.Run[target.Portal],
) *Runner[T] {
	return &Runner[T]{runners: runners}
}

func (r Runner[T]) Run(ctx context.Context, portal T) (err error) {
	for _, runner := range r.runners {
		err = runner(ctx, portal)
		if !errors.Is(err, target.ErrNotTarget) {
			return
		}
	}
	return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(portal), portal.Path())
}
