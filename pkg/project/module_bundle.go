package project

import (
	"archive/zip"
	"bytes"
	"github.com/cryptopunkscc/go-astral-js/pkg/bundle"
	"io/fs"
	"log"
	"path"
)

type Bundle struct {
	*Module
	manifest bundle.Manifest
}

func (m *Module) Bundle() (b *Bundle, err error) {

	file, err := fs.ReadFile(m.files, path.Base(m.Path()))
	if err != nil {
		return
	}

	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		log.Println("reader err", err, m.Path())
		return
	}
	manifest, err := bundle.ReadManifestFs(reader)
	if err != nil {
		return
	}
	b = &Bundle{
		Module:   newModuleFS(m.Path(), reader),
		manifest: manifest,
	}
	return
}

func (b *Bundle) App() {}

func (b *Bundle) Manifest() bundle.Manifest {
	return b.manifest
}
