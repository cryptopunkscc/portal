package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/resolve/source"
	"io/fs"
)

type pkg struct {
	target.Source
	pkg target.Source
}

func (s *pkg) Package() target.Source { return s.pkg }

func Resolve(src target.Source) (t target.Bundle, err error) {
	if src.IsDir() {
		return nil, errors.New("not a file")
	}
	zipReader, err := zipFromSource(src)
	if err != nil {
		return
	}
	s := source.FS(zipReader, src.Abs())
	t = &pkg{
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
