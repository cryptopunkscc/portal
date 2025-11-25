package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/fs"

	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target/source"
)

func Resolve(src target.Source) (t target.Bundle, err error) {
	defer plog.TraceErr(&err)
	if src.IsDir() {
		return nil, errors.New("not a file")
	}
	reader, err := Reader(src)
	if err != nil {
		return
	}
	unpacked := source.FS(reader, src.Abs())
	t = &File_{
		Source: unpacked,
		file:   src,
	}
	return

}

func Reader(src target.Source) (r *zip.Reader, err error) {
	var file []byte
	if file, err = fs.ReadFile(src.FS(), src.Path()); err != nil {
		return
	}
	readerAt := bytes.NewReader(file)
	size := int64(len(file))
	return zip.NewReader(readerAt, size)
}
