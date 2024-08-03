package target

import (
	"io/fs"
	"path/filepath"
)

func Any[T Source](of ...func(Source) (Source, error)) Resolve[T] {
	return Combine[Source, T](of...)
}

func Try[A Source, B Source](f func(A) (B, error)) Resolve[Source] {
	return func(arg Source) (s Source, err error) {
		a, ok := arg.(A)
		if !ok {
			return
		}
		return f(a)
	}
}

func Skip(names ...string) func(source Source) (result Source, err error) {
	return func(source Source) (result Source, err error) {
		for _, n := range names {
			if filepath.Base(source.Abs()) == n {
				return nil, fs.SkipDir
			}
		}
		return
	}
}
