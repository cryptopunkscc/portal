package project

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"log"
)

var _ target.App = &Bundle{}

type Bundle struct {
	target.Source
	manifest *bundle.Manifest
}

func NewBundle(abs string) (b *Bundle, err error) {
	return ResolveBundle(NewModule(abs))
}

func ResolveBundle(source target.Source) (b *Bundle, err error) {
	if !source.Type().Is(target.Bundle) {
		err = errors.New("not a bundle")
		return
	}

	file, err := fs.ReadFile(source.Files(), source.Path())
	if err != nil {
		return
	}

	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		log.Println("reader err", err, source.Path())
		return
	}
	s := NewModuleFS(reader, source.Path())
	s.abs = source.Abs()
	m, err := bundle.ReadManifestFs(reader)
	if err != nil {
		return
	}
	b = &Bundle{
		Source:   s,
		manifest: &m,
	}
	return
}

func (b *Bundle) App() {}

func (b *Bundle) Manifest() *bundle.Manifest {
	return b.manifest
}
