package source

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

func Collect(src Source, constructors ...Constructor) (out List[Source]) {
	return CollectT[Source](src, constructors...)
}

func CollectIt[T Constructor](src Source, constructor T) (out List[T]) {
	return CollectT[T](src, constructor)
}

func CollectT[T Source](src Source, constructors ...Constructor) (out List[T]) {
	if src == nil {
		return
	}

	if src != src.Ref_() {
		if o, err := ResolveT[T](src, constructors...); err == nil {
			out = append(out, o)
		}
		return
	}

	ref := *src.Ref_()
	if ref.Fs == nil {
		return
	}
	if ref.Path == "" {
		ref.Path = "."
	}
	_ = afero.Walk(ref.Fs, ref.Path, func(p string, _ fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(p, "node_modules") {
			return fs.SkipDir
		}
		ref.Path = p
		source, err := ResolveT[T](&ref, constructors...)
		if isSkip(err) {
			return err
		}
		if err == nil {
			out = append(out, source)
		}
		return nil
	})
	return
}

func ResolveIt[T Constructor](src Source, constructor T) (out T, err error) {
	return ResolveT[T](src, constructor)
}

func Resolve(src Source, constructors ...Constructor) (out Source, err error) {
	return ResolveT[Source](src, constructors...)
}

func ResolveT[T Source](src Source, constructors ...Constructor) (out T, err error) {
	if len(constructors) == 0 {
		return out, fs.ErrInvalid
	}
	for _, constructor := range constructors {
		if out, err = Cast[T](constructor.New()); err == nil {
			if err = out.ReadSrc(src); err == nil {
				return
			}
		}
		if isSkip(err) {
			return
		}
	}
	return out, fs.ErrInvalid
}

func isSkip(err error) bool {
	return errors.Is(err, fs.SkipDir) || errors.Is(err, fs.SkipAll)
}

func Cast[T Source](src Source) (out T, err error) {
	out, ok := src.(T)
	if !ok {
		return out, fs.ErrInvalid
	}
	return
}
