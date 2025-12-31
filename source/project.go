package source

import (
	"encoding/json"
	"path"
	"path/filepath"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
	"go.nhat.io/aferocopy/v2"
	"gopkg.in/yaml.v3"
)

type Project struct {
	Ref
	Metadata ProjectMetadata
}

func (p Project) Dist() (a App, err error) {
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

func (p *Project) ReadSrc(src Source) (err error) {
	return Readers{&p.Ref, &p.Metadata}.ReadSrc(src)
}

func (p *Project) WriteRef(ref Ref) (err error) {
	return Writers{&p.Metadata, &p.Ref}.WriteRef(ref)
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

func (m *ProjectMetadata) ReadSrc(src Source) (err error) {
	defer plog.TraceErr(&err)
	ref := *src.Ref_()
	bytes, err := afero.ReadFile(ref.Fs, path.Join(ref.Path, "dev.portal.json"))
	if err == nil {
		return json.Unmarshal(bytes, m)
	}
	bytes, err = afero.ReadFile(ref.Fs, path.Join(ref.Path, "dev.portal.yaml"))
	if err == nil {
		return yaml.Unmarshal(bytes, m)
	}
	bytes, err = afero.ReadFile(ref.Fs, path.Join(ref.Path, "dev.portal.yml"))
	if err == nil {
		return yaml.Unmarshal(bytes, m)
	}
	return
}

func (m *ProjectMetadata) WriteRef(ref Ref) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	return afero.WriteFile(ref.Fs, path.Join(ref.Path, "dev.portal.json"), bytes, 0644)
}
