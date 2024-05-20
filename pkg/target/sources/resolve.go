package sources

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/dist"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/npm"
	"github.com/cryptopunkscc/go-astral-js/pkg/target/project"
	"io/fs"
)

func FromPath[T target.Source](src string) (in <-chan T) {
	return target.Stream[T](resolve[T](), target.NewModule(src))
}

func FromFS[T target.Source](src fs.FS) (in <-chan T) {
	return target.Stream[T](resolve[T](), target.NewModuleFS(src))
}

func resolve[T target.Source]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Lift(target.Try(npm.ResolveNodeModule))(
			target.Try(project.Resolve)),
		target.Try(dist.Resolve),
	)
}
