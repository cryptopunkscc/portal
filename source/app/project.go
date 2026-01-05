package app

import (
	"path"
	"path/filepath"

	"github.com/cryptopunkscc/portal/source"
	"go.nhat.io/aferocopy/v2"
)

type Project struct {
	source.Ref
	Metadata ProjectMetadata
}

func (p Project) Dist() (a Dist, err error) {
	err = a.ReadSrc(p.Sub("dist"))
	return
}

func (p Project) Pack() (err error) {
	app, err := p.Dist()
	if err != nil {
		return
	}
	return app.Bundle().WriteRef(*p.Sub("build"))
}

func (p *Project) ReadSrc(src source.Source) (err error) {
	return source.Readers{&p.Ref, &p.Metadata}.ReadSrc(src)
}

func (p *Project) WriteRef(ref source.Ref) (err error) {
	return source.Writers{&p.Metadata, &p.Ref}.WriteRef(ref)
}

func (p Project) Build(dstPath string) (err error) {
	// copy app icon into dist
	if len(p.Metadata.Icon) > 0 {
		iconName := "icon" + filepath.Ext(p.Metadata.Icon)
		if err = aferocopy.Copy(
			path.Join(p.Path, p.Metadata.Icon),
			path.Join(p.Path, dstPath, iconName),
			aferocopy.Options{SrcFs: p.Fs, DestFs: p.Fs},
		); err != nil {
			return
		}
		p.Metadata.Icon = iconName
	}

	// write portal.json into dist
	if err = p.Metadata.Manifest.WriteRef(*p.Sub(dstPath)); err != nil {
		return
	}

	return
}

type ProjectMetadata struct {
	Metadata `json:",inline" yaml:",inline"`
	Builds   InnerBuilds `yaml:"build,omitempty" json:"build,omitempty"`
}

func (m *ProjectMetadata) ReadSrc(src source.Source) (err error) {
	return metadataReadSrc(m, "dev.portal", src)
}

func (m *ProjectMetadata) WriteRef(ref source.Ref) (err error) {
	return metadataWriteRef(m, "dev.portal", ref)
}
