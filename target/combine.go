package target

import (
	"errors"
	"io/fs"
)

func Combine[A any, T any](of ...func(A) (A, error)) func(A) (T, error) {
	return func(entry A) (s T, err error) {
		for _, f := range of {
			var v A
			v, err = f(entry)
			if err != nil {
				if errors.Is(err, fs.SkipDir) {
					return
				}
				err = nil
				continue
			}
			ok := false
			if s, ok = any(v).(T); ok {
				return
			}
		}
		return
	}
}
