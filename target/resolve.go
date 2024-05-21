package target

import (
	"errors"
	"io/fs"
	"path"
)

func Any[T Source](of ...func(Source) (Source, error)) Resolve[T] {
	return func(entry Source) (s T, err error) {
		for _, f := range of {
			var v Source
			v, err = f(entry)
			if err != nil {
				if errors.Is(err, fs.SkipDir) {
					return
				}
				err = nil
				continue
			}
			ok := false
			if s, ok = v.(T); ok {
				return
			}
		}
		return
	}
}

func Lift(
	from func(Source) (Source, error),
) func(to ...func(Source) (Source, error)) func(Source) (Source, error) {
	return func(to ...func(Source) (Source, error)) func(Source) (Source, error) {
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
}

func Try[A Source, B Source](f func(A) (B, error)) func(Source) (Source, error) {
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
			if path.Base(source.Path()) == n {
				return nil, fs.SkipDir
			}
		}
		return
	}
}
