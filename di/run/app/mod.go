package app

import (
	"github.com/cryptopunkscc/portal/di/run"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target2/apps"
)

type Module[T App_] struct{ run.Module[T] }

func (d *Module[T]) TargetResolve() Resolve[T] { return apps.Resolver[T]() }
