package app

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/apps"
)

type Module[T App_] struct{}

func (d *Module[T]) Resolver() Resolve[T] { return apps.Resolver[T]() }
