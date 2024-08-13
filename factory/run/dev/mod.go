package dev

import (
	"github.com/cryptopunkscc/portal/resolve/sources"
	. "github.com/cryptopunkscc/portal/target"
)

type Module[T Portal_] struct{}

func (d *Module[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }
