package dev

import (
	"github.com/cryptopunkscc/portal/di/run"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/portals"
)

type Module[T Portal] struct{ run.Module[T] }

func (d *Module[T]) TargetFinder() Finder[T] { return portals.NewFind[T] }
