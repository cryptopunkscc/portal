package dev

import (
	"github.com/cryptopunkscc/portal/factory/run"
	"github.com/cryptopunkscc/portal/resolve/sources"
	. "github.com/cryptopunkscc/portal/target"
)

type Module[T Portal_] struct{ run.Module[T] }

func (d *Module[T]) TargetResolve() Resolve[T] { return sources.Resolver[T]() }
