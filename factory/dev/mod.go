package dev

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/sources"
)

type Module[T Portal_] struct{}

func (d *Module[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }
