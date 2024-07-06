package dev

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/portal/target"
	"reflect"
)

type Runner struct {
	runners []target.Run[target.Portal]
}

func NewRunner(
	runners ...target.Run[target.Portal],
) *Runner {
	return &Runner{runners: runners}
}

func (r Runner) Run(ctx context.Context, portal target.Portal) (err error) {
	for _, runner := range r.runners {
		err = runner(ctx, portal)
		if !errors.Is(err, target.ErrNotTarget) {
			return
		}
	}
	return fmt.Errorf("invalid target %v: %v", reflect.TypeOf(portal), portal.Path())
}
