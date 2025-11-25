package target

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"

	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Resolve[T any] func(src Source) (result T, err error)

var Resolve_ Resolve[any] = func(Source) (result any, err error) { return }

func (resolve Resolve[T]) Resolve(src Source) (result T, err error) { return resolve(src) }

// List all Source from a given dir.
func (resolve Resolve[T]) List(from ...Source) (out []T) {
	for _, src := range from {
		o := resolve.list(src)
		out = append(out, o...)
	}
	return
}

// list all Source from a given dir.
func (resolve Resolve[T]) list(from Source) (out []T) {
	if !from.IsDir() {
		if t, err := resolve(from); err == nil {
			return append(out, t)
		}
	}
	_ = fs.WalkDir(from.FS(), ".", func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println("Resolve list", err)
			return err
		}
		m, err := from.Sub(src)
		if err != nil {
			log.Println("Resolve sub", err)
			return nil
		}
		s, err := resolve(m)
		if any(s) != nil && !reflect.ValueOf(s).IsZero() {
			out = append(out, s)
		}
		if errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll) {
			return errors.Unwrap(err)
		}
		return nil
	})
	return
}

func (resolve Resolve[T]) Try(src Source) (result Source, err error) {
	r, err := resolve(src)
	if err != nil {
		return
	}
	result, ok := any(r).(Source)
	if !ok {
		err = ErrNotTarget
	}
	return
}

func Any[T Source](of ...Resolve[Source]) Resolve[T] {
	return func(entry Source) (s T, err error) {
		defer plog.TraceErr(&err)
		var ok bool
		if s, ok = any(err).(T); ok {
			return
		}

		for _, f := range of {
			var v Source
			if v, err = f(entry); err != nil {
				if errors.Is(err, fs.SkipDir) {
					return
				}
				err = nil
				continue
			}
			if s, ok = any(v).(T); ok {
				return
			}
		}
		err = ErrNotTarget
		return
	}
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

func Skip(names ...string) func(Source) (Source, error) {
	return func(source Source) (result Source, err error) {
		for _, n := range names {
			if filepath.Base(source.Abs()) == n {
				return nil, fs.SkipDir
			}
		}
		return
	}
}
