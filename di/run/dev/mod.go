package dev

import (
	"github.com/cryptopunkscc/portal/di/run"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/sources"
)

type Module[T Base] struct{ run.Module[T] }

func (d *Module[T]) TargetResolve() Resolve[T] { return sources.Resolver[T]() }
