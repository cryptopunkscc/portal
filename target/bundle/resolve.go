package bundle

import (
	"archive/zip"
	"bytes"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
	"log"
)

func New(abs string) (b target.Bundle, err error) {
	return Resolve(targetSource.New(abs))
}

func Resolve(t target.Source) (b target.Bundle, err error) {
	t = t.Lift()
	if !t.Type().Is(target.TypeBundle) {
		return nil, target.ErrNotTarget
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
	s := targetSource.Resolve(reader, t.Path(), t.Abs())
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
