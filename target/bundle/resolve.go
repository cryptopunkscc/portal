package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/manifest"
	targetSource "github.com/cryptopunkscc/portal/target/source"
	"io/fs"
)

var ErrNotBundle = errors.New("not a bundle")

func Resolve(src target.Source) (bundle target.Bundle, err error) {
	src = src.Lift()
	if !src.Type().Is(target.TypeBundle) {
		return nil, ErrNotBundle
	}

	var zipFs fs.FS

	// FIXME
	//if path.IsAbs(src.Abs()) {
	//	if zipFs, err = zip.OpenReader(src.Abs()); err != nil {
	//		return
	//	}
	//} else {
	if zipFs, err = zipFromSource(src); err != nil {
		return
	}
	//}

	s := targetSource.FromFS(zipFs, src.Path(), src.Abs())
	m, err := manifest.Read(zipFs)
	if err != nil {
		return
	}
	bundle = &source{Source: s, manifest: &m}
	switch {
	case bundle.Type().Is(target.TypeFrontend):
		bundle = &frontend{Bundle: bundle}
	case bundle.Type().Is(target.TypeBackend):
		bundle = &backend{Bundle: bundle}
	default:
		e := &executable{Bundle: bundle}
		if e.Exec, err = exec.Resolve(bundle); err != nil {
			return
		}
		bundle = e
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
