package resolve

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
)

func FromPath[T target.Source](src string) (in <-chan T) {
	return target.Stream[T](Source, target.NewModule(src))
}

func FromFS[T target.Source](src fs.FS) (in <-chan T) {
	return target.Stream[T](Source, target.NewModuleFS(src))
}
