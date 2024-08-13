package app

import (
	"github.com/cryptopunkscc/portal/resolve/apps"
	. "github.com/cryptopunkscc/portal/target"
)

type Module[T App_] struct{}

func (d *Module[T]) Resolver() Resolve[T] { return apps.Resolver[T]() }
