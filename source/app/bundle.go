package app

import (
	"bytes"
	"io"
	"path"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/spf13/afero"
)

type Bundle struct {
	Dist
	source.Zip
}

func (b Bundle) New() (src source.Source) {
	return &b
}

func init() {
	_ = astral.DefaultBlueprints.Add(&Bundle{})
}

func (b *Bundle) ReadSrc(src source.Source) (err error) {
	if err = b.Zip.ReadSrc(src); err != nil {
		return
	}
	if err = b.Dist.ReadSrc(&source.Ref{Fs: b.Unpacked}); err != nil {
		return
	}
	return
}

func (b *Bundle) WriteRef(ref source.Ref) (err error) {
	if b.Dist.Fs == nil {
		b.Dist.Fs = afero.NewMemMapFs()
	}
	if err = b.Dist.WriteRef(b.Dist.Ref); err != nil {
		return
	}

	b.Zip.Unpacked = b.Dist.Fs
	if len(b.Dist.Path) > 0 {
		b.Zip.Unpacked = afero.NewBasePathFs(b.Fs, b.Path)
	}

	ref.Path = path.Join(ref.Path, b.BundleName())
	return b.Zip.WriteRef(ref)
}

func (b *Bundle) ObjectType() string { return "app.bundle" }

func (b *Bundle) WriteTo(w io.Writer) (n int64, err error) {
	defer plog.TraceErr(&err)
	i, err := w.Write(b.Blob)
	return int64(i), err
}

func (b *Bundle) ReadFrom(r io.Reader) (n int64, err error) {
	defer plog.TraceErr(&err)
	blob, err := io.ReadAll(r)
	if err != nil {
		return
	}
	return b.Zip.ReadFrom(bytes.NewReader(blob))
}

func (b Bundle) Publish(objects *astrald.ObjectsClient) (info ReleaseInfo, err error) {
	release := ReleaseMetadata{
		Release: b.Metadata.Release,
		Target:  b.Metadata.Target,
	}
	if release.BundleID, err = source.ObjectsCommit(objects, bytes.NewReader(b.Blob)); err != nil {
		return
	}
	if release.ManifestID, err = source.ObjectsCommit(objects, &b.Metadata.Manifest); err != nil {
		return
	}
	if info.ReleaseID, err = source.ObjectsCommit(objects, &release); err != nil {
		return
	}

	info.Manifest = b.Metadata.Manifest
	info.ReleaseMetadata = release
	return
}
