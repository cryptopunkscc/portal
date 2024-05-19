package portal

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"io/fs"
	"log"
)

type Bundle struct {
	target.Source
	manifest *target.Manifest
}

var _ target.Bundle = (*Bundle)(nil)

type FrontendBundle struct {
	target.Frontend
	target.Bundle
}

type BackendBundle struct {
	target.Backend
	target.Bundle
}

func NewBundle(abs string) (b target.Bundle, err error) {
	return ResolveBundle(target.NewModule(abs))
}

func ResolveBundle(source target.Source) (b target.Bundle, err error) {
	source = source.Lift()
	if !source.Type().Is(target.TypeBundle) {
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
	s := target.NewModuleFS(reader, source.Path(), source.Abs())
	m, err := target.ReadManifestFs(reader)
	if err != nil {
		return
	}
	b = &Bundle{
		Source:   s,
		manifest: &m,
	}
	switch {
	case b.Type().Is(target.TypeFrontend):
		b = &FrontendBundle{Bundle: b}
	case b.Type().Is(target.TypeBackend):
		b = &BackendBundle{Bundle: b}
	}
	return
}

func (b *Bundle) IsApp() {}

func (b *Bundle) IsBundle() {}

func (b *Bundle) Manifest() *target.Manifest {
	return b.manifest
}
