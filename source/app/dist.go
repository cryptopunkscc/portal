package app

import (
	"fmt"

	"github.com/cryptopunkscc/portal/source"
	"github.com/spf13/afero"
)

type Dist struct {
	source.Ref
	Metadata
}

func (a Dist) New() (src source.Source) {
	return &a
}

func (a Dist) Dist() Dist {
	return a
}

func (a Dist) Bundle() *Bundle {
	return &Bundle{Dist: a, Zip: source.Zip{
		Unpacked: afero.NewBasePathFs(a.Fs, a.Path)}}
}

func (a *Dist) BundleName() string {
	t := a.Metadata.Target
	platform := ""
	if len(t.OS) > 0 {
		platform += "_" + t.OS
		if len(t.Arch) > 0 {
			platform += "_" + t.Arch
		}
	}

	m := a.Metadata
	version := fmt.Sprintf("%d.%d.%d", m.Version, m.Api.Version, m.Release.Version)

	return fmt.Sprintf("%s_%s%s.portal", m.Package, version, platform)
}

func (a *Dist) ReadSrc(src source.Source) (err error) {
	return source.Readers{&a.Ref, &a.Metadata}.ReadSrc(src)
}

func (a *Dist) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&a.Ref, &a.Metadata}.WriteRef(ref)
}

func (a *Dist) ReadFs(fS afero.Fs) (err error) {
	return source.Readers{&a.Ref, &a.Metadata}.ReadSrc(&source.Ref{Fs: fS})
}
