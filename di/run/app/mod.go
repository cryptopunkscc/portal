package app

import (
	"github.com/cryptopunkscc/portal/di/run"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apps"
)

type Module[T App] struct{ run.Module[T] }

func (d *Module[T]) TargetFinder() Finder[T] { return apps.NewFind[T] }
