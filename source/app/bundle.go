package app

import (
	"bytes"
	"context"
	"io"
	"path"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/cryptopunkscc/portal/core/apphost"
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
	_ = astral.Add(&Bundle{})
}

func (b *Bundle) ReadSrc(src source.Source) (err error) {
	if err = b.Zip.ReadSrc(src); err != nil {
		return
	}
	if err = b.Dist.ReadFs(b.Unpacked); err != nil {
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
	if n, err = b.Zip.ReadFrom(bytes.NewReader(blob)); err != nil {
		return
	}
	if err = b.Dist.ReadFs(b.Unpacked); err != nil {
		return
	}
	return
}

func (b Bundle) Publish(ctx context.Context, objects *apphost.ObjectsClient) (info ReleaseInfo, err error) {
	release := ReleaseMetadata{
		Release: b.Metadata.Release,
		Target:  b.Metadata.Target,
	}
	repo := "main"
	actx := astral.NewContext(ctx)
	if release.BundleID, err = objects.Store(actx, repo, &b); err != nil {
		return
	}
	if release.ManifestID, err = objects.Store(actx, repo, &b.Metadata.Manifest); err != nil {
		return
	}
	if info.ReleaseID, err = objects.Store(actx, repo, &release); err != nil {
		return
	}

	info.Manifest = b.Metadata.Manifest
	info.ReleaseMetadata = release
	return
}

// Fixme: why objects.Module.OpCreate gets stuck on ch.Send(&astral.Ack{})?
func (b Bundle) Publish2(ctx context.Context, objects *objects.Client) (info ReleaseInfo, err error) {
	release := ReleaseMetadata{
		Release: b.Metadata.Release,
		Target:  b.Metadata.Target,
	}
	actx := astral.NewContext(ctx)
	writer, err := objects.Create(actx, "main", 0)
	if err != nil {
		return
	}

	if _, err = b.WriteTo(writer); err != nil {
		return
	}
	if release.BundleID, err = writer.Commit(); err != nil {
		return
	}

	if _, err = b.Manifest.WriteTo(writer); err != nil {
		return
	}
	if release.ManifestID, err = writer.Commit(); err != nil {
		return
	}

	if _, err = release.WriteTo(writer); err != nil {
		return
	}
	if info.ReleaseID, err = writer.Commit(); err != nil {
		return
	}

	info.Manifest = b.Metadata.Manifest
	info.ReleaseMetadata = release
	return
}
