package portal

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/app"
	"github.com/cryptopunkscc/portal/resolve/project"
)

var Resolve_ = target.Any[target.Portal_](
	Resolver(target.Resolve_).Try,
)

func Resolver[T any](resolveT target.Resolve[T]) target.Resolve[target.Portal[T]] {
	return target.Any[target.Portal[T]](
		app.Resolver[T](resolveT).Try,
		project.Resolver(resolveT).Try,
	)
}
