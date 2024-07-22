package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

type of[T any] struct {
	target.Dist[T]
	bundle target.Bundle
}

func (t of[T]) Package() target.Source { return t.bundle.Package() }

func Resolver[T any](resolve target.Resolve[target.Dist[T]]) target.Resolve[target.AppBundle[T]] {
	return func(src target.Source) (app target.AppBundle[T], err error) {
		b, err := Resolve(src)
		if err != nil {
			return
		}
		td := &of[T]{}
		if td.Dist, err = resolve(b); err != nil {
			return
		}
		td.bundle = b
		app = td
		return
	}
}

type raw struct {
	target.Source
	pkg target.Source
}

func (s *raw) Package() target.Source { return s.pkg }

func Resolve(src target.Source) (t target.Bundle, err error) {
	if src.IsDir() {
		return nil, errors.New("not a file")
	}
	zipReader, err := zipFromSource(src)
	if err != nil {
		return
	}
	s := source.FS(zipReader, src.Abs())
	t = &raw{
		Source: s,
		pkg:    src,
	}
	return

}

func zipFromSource(src target.Source) (r *zip.Reader, err error) {
	var file []byte
	if file, err = fs.ReadFile(src.Files(), src.Path()); err != nil {
		return
	}
	readerAt := bytes.NewReader(file)
	size := int64(len(file))
	return zip.NewReader(readerAt, size)
}
