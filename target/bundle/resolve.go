package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
	"log"
)

func FromPath(src string) (b target.Bundle, err error) {
	return Resolve(targetSource.FromPath(src))
}

var ErrNotBundle = errors.New("not a bundle")

func Resolve(t target.Source) (b target.Bundle, err error) {
	t = t.Lift()
	if !t.Type().Is(target.TypeBundle) {
		return nil, ErrNotBundle
	}

	file, err := fs.ReadFile(t.Files(), t.Path())
	if err != nil {
		return
	}

	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		log.Println("reader err", err, t.Path())
		return
	}
	s := targetSource.FromFS(reader, t.Path(), t.Abs())
	m, err := manifest.Read(reader)
	if err != nil {
		return
	}
	b = &source{
		Source:   s,
		manifest: &m,
	}
	switch {
	case b.Type().Is(target.TypeFrontend):
		b = &frontend{Bundle: b}
	case b.Type().Is(target.TypeBackend):
		b = &backend{Bundle: b}
	}
	return
}
