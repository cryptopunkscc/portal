package source

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"path"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type Zip struct {
	Unpacked afero.Fs
	Blob     []byte
	File     Ref
	ObjectID *astral.ObjectID
}

func (z *Zip) ReadSrc(src Source) (err error) {
	defer plog.TraceErr(&err)
	if err = z.File.ReadSrc(src); err != nil {
		return
	}

	blob, err := afero.ReadFile(z.File.Fs, z.File.Path)
	if err != nil {
		return
	}

	_, err = z.ReadFrom(bytes.NewReader(blob))
	return
}

func (z *Zip) ReadFrom(r io.Reader) (n int64, err error) {
	blob, err := io.ReadAll(r)
	reader, err := zip.NewReader(bytes.NewReader(blob), int64(len(blob)))
	if err != nil {
		return
	}
	z.Unpacked = afero.FromIOFS{FS: reader}
	z.Blob = blob
	return
}

func (z *Zip) WriteRef(ref Ref) (err error) {
	defer plog.TraceErr(&err)
	if z.Unpacked == nil {
		return errors.New("zip.WriteRef: Unpacked is nil")
	}

	buffer := bytes.Buffer{}
	zipWriter := zip.NewWriter(&buffer)
	if err = zipWriter.AddFS(afero.NewIOFS(z.Unpacked)); err != nil {
		return
	}
	if err = zipWriter.Close(); err != nil {
		return
	}

	if err = ref.Fs.MkdirAll(path.Dir(ref.Path), 0755); err != nil {
		return
	}
	if err = afero.WriteFile(ref.Fs, ref.Path, buffer.Bytes(), 0644); err != nil {
		return
	}

	z.File = ref
	return
}

func (z Zip) Publish(objects *astrald.ObjectsClient) (err error) {
	writer, err := objects.Create("", len(z.Blob))
	if err != nil {
		return
	}
	if _, err = writer.Write(z.Blob); err != nil {
		return
	}
	z.ObjectID, err = writer.Commit()
	return
}
