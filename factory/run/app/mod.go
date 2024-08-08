package app

import (
	"github.com/cryptopunkscc/portal/factory/run"
	"github.com/cryptopunkscc/portal/resolve/apps"
	. "github.com/cryptopunkscc/portal/target"
)

type Module[T App_] struct{ run.Module[T] }

func (d *Module[T]) TargetResolve() Resolve[T] { return apps.Resolver[T]() }
