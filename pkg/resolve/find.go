package resolve

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/project"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

func FromPath[T target.Source](src string) (in <-chan T) {
	return target.Stream[T](Dev[T](), target.NewModule(src))
}

func FromFS[T target.Source](src fs.FS) (in <-chan T) {
	return target.Stream[T](Dev[T](), target.NewModuleFS(src))
}

func Dev[T target.Source]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Try(project.SkipNodeModulesDir),
		target.Try(portal.ResolveBundle),
		target.Lift(target.Try(project.ResolveNodeModule))(
			target.Try(project.ResolvePortalModule)),
		target.Try(portal.ResolveDist),
	)
}

func App[T target.Source]() func(target.Source) (T, error) {
	return target.Any[T](
		target.Try(project.SkipNodeModulesDir),
		target.Try(portal.ResolveBundle),
		target.Try(portal.ResolveDist),
	)
}
