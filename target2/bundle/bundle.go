package bundle

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/cryptopunkscc/portal/target2"
	"github.com/cryptopunkscc/portal/target2/source"
	"io/fs"
)

type bundle struct{ target2.Source }

func (s *bundle) IsBundle() {}
func (s *bundle) IsApp()    {}

func Resolve(src target2.Source) (t target2.Bundle, err error) {
	if src.IsDir() {
		return nil, errors.New("not a file")
	}
	r, err := zipFromSource(src)
	if err != nil {
		return
	}
	s := source.FS(r, src.Abs())
	t = &bundle{s}
	return

}

func zipFromSource(src target2.Source) (r *zip.Reader, err error) {
	var file []byte
	if file, err = fs.ReadFile(src.Files(), src.Path()); err != nil {
		return
	}
	readerAt := bytes.NewReader(file)
	size := int64(len(file))
	return zip.NewReader(readerAt, size)
}
