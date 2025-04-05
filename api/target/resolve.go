package target

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
)

type Resolve[T any] func(src Source) (result T, err error)

type Resolver[T any] interface {
	Resolve(src Source) (result T, err error)
}

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
	_ = fs.WalkDir(from.FS(), from.Path(), func(src string, d fs.DirEntry, err error) error {
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
		if errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll) {
			return err
		}
		if any(s) != nil && !reflect.ValueOf(s).IsNil() {
			out = append(out, s)
		}
		return nil
	})
	return
}

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
