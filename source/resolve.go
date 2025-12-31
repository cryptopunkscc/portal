package source

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

func Collect(src Source, types ...Type) (out List[Source]) {
	return CollectT[Source](src, types...)
}

func CollectIt[T Type](src Source, t T) (out List[T]) {
	return CollectT[T](src, t)
}

func CollectT[T Source](src Source, types ...Type) (out List[T]) {
	if src == nil {
		return
	}

	if src != src.Ref_() {
		if o, err := ResolveT[T](src, types...); err == nil {
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
		source, err := ResolveT[T](&ref, types...)
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

func ResolveIt[T Type](src Source, t T) (out T, err error) {
	return ResolveT[T](src, t)
}

func Resolve(src Source, types ...Type) (out Source, err error) {
	return ResolveT[Source](src, types...)
}

func ResolveT[T Source](src Source, types ...Type) (out T, err error) {
	if len(types) == 0 {
		return out, fs.ErrInvalid
	}
	for _, constructor := range types {
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
