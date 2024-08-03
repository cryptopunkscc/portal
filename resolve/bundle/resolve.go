package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/target"
	"io/fs"
)

type of[T any] struct{ target.Dist[T] }

func (t of[T]) IsBundle() {}

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
		app = td
		return
	}
}

type raw struct{ target.Source }

func (s *raw) IsBundle() {}

func Resolve(src target.Source) (t target.Bundle, err error) {
	if src.IsDir() {
		return nil, errors.New("not a file")
	}
	r, err := zipFromSource(src)
	if err != nil {
		return
	}
	s := source.FS(r, src.Abs())
	t = &raw{s}
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
