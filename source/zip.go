package source

import (
	"archive/zip"
	"bytes"
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type ZipBundle struct {
	Source
	ZipFs afero.Fs
}

func (b *ZipBundle) ReadSource(source Source) (err error) {
	b.Source = source
	return b.ReadFs(source.Fs)
}

func (b *ZipBundle) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	if err = b.Source.ReadFs(files); err != nil {
		return
	}
	if err = b.SetZipReader(); err != nil {
		return
	}
	return
}

func (b *ZipBundle) WriteZipFs(out afero.Fs) (err error) {
	if err = b.Source.WriteZipFs(out); err != nil {
		return
	}
	if err = b.SetZipReader(); err != nil {
		return
	}
	return
}

func (b *ZipBundle) SetZipReader() (err error) {
	defer plog.TraceErr(&err)
	var file []byte
	if file, err = fs.ReadFile(afero.IOFS{Fs: b.Fs}, b.Name); err != nil {
		return
	}
	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))
	if err != nil {
		return
	}
	b.ZipFs = afero.NewCopyOnWriteFs(afero.FromIOFS{FS: reader}, afero.NewMemMapFs())
	return
}
