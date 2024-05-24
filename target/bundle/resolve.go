package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/manifest"
	targetSource "github.com/cryptopunkscc/go-astral-js/target/source"
	"io/fs"
)

func FromPath(src string) (bundle target.Bundle, err error) {
	return Resolve(targetSource.FromPath(src))
}

var ErrNotBundle = errors.New("not a bundle")

func Resolve(src target.Source) (bundle target.Bundle, err error) {
	src = src.Lift()
	if !src.Type().Is(target.TypeBundle) {
		return nil, ErrNotBundle
	}

	file, err := fs.ReadFile(src.Files(), src.Path())
	if err != nil {
		return
	}

	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		return
	}
	s := targetSource.FromFS(reader, src.Path(), src.Abs())
	m, err := manifest.Read(reader)
	if err != nil {
		return
	}
	bundle = &source{
		Source:   s,
		manifest: &m,
	}
	switch {
	case bundle.Type().Is(target.TypeFrontend):
		bundle = &frontend{Bundle: bundle}
	case bundle.Type().Is(target.TypeBackend):
		bundle = &backend{Bundle: bundle}
	}
	return
}
