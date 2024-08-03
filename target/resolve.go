package target

import (
	"io/fs"
	"log"
	"path/filepath"
)

func Any[T Source](of ...func(Source) (Source, error)) Resolve[T] {
	return Combine[Source, T](of...)
}

type ResolveSource Resolve[Source]

func (from ResolveSource) Lift(to ...func(Source) (Source, error)) func(Source) (Source, error) {
	return func(src Source) (s Source, err error) {
		if s, err = from(src); err != nil {
			return
		}
		ss := s
		for _, o := range to {
			if s, err = o(ss); err == nil {
				return
			}
		}
		err = nil
		s = ss
		return
	}
}

func Try[A Source, B Source](f func(A) (B, error)) ResolveSource {
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
				log.Println("skip node module", filepath.Base(source.Abs()), source.Abs())
				return nil, fs.SkipDir
			}
		}
		return
	}
}
