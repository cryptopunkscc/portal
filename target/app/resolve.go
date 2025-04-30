package app

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/dist"
)

var Resolve_ = Resolver(target.Resolve_)

func Resolver[T any](resolveType target.Resolve[T]) target.Resolve[target.App[T]] {
	resolveDist := dist.Resolver[T](resolveType)
	resolveBundle := bundle.Resolver[T](resolveDist)
	return target.Any[target.App[T]](
		target.Try(resolveDist),
		target.Try(resolveBundle),
	)
}
