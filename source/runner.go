package source

import (
	"io/fs"
)

type Runner[T Source] struct {
	Provider Provider
	Types    Types
	Handler  func([]T) error
}

func (r Runner[T]) Run(src string) error {
	if r.Provider == nil {
		r.Provider = OsFs
	}
	s := r.Provider.GetSource(src)
	if s == nil {
		return fs.ErrNotExist
	}
	sources := CollectT[T](s, r.Types...)

	return r.Handler(sources)
}
