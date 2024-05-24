package apps

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
)

func Resolve[T target.App]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Try(dist.Resolve),
	)
}
