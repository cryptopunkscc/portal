package sources

import (
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/bundle"
	"github.com/cryptopunkscc/go-astral-js/target/dist"
	"github.com/cryptopunkscc/go-astral-js/target/npm"
	"github.com/cryptopunkscc/go-astral-js/target/project"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

func FromPath[T target.Source](src string) []T {
	return source.List[T](Resolve[T](), source.FromPath(src))
}

func FromFS[T target.Source](src fs.FS) []T {
	return source.List[T](Resolve[T](), source.FromFS(src))
}

func Resolve[T target.Source]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Skip("node_modules"),
		target.Try(bundle.Resolve),
		target.Lift(target.Try(npm.ResolveNodeModule))(
			target.Try(project.Resolve)),
		target.Try(dist.Resolve),
	)
}