package apps

import (
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/dist"
)

func Resolve[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(dist.Resolve),
	)
}
