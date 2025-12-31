package source

import (
	"encoding/json"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Project struct {
	Source

	Manifest ProjectMetadata
}

func (p *Project) ReadFs(files afero.Fs) (err error) {
	return FSReaders{&p.Source, &p.Manifest}.ReadFs(files)
}

func (p *Project) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&p.Source, &p.Manifest}.WriteFs(dir)
}

func (p *Project) WriteOS(dir string) (err error) {
	return p.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type ProjectMetadata struct {
	Manifest `json:",inline" yaml:",inline"`
	Builds   InnerBuilds `yaml:"build,omitempty" json:"build,omitempty"`
}

func (m *ProjectMetadata) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)

	bytes, err := afero.ReadFile(files, "dev.portal.json")
	if err == nil {
		return json.Unmarshal(bytes, m)
	}

	bytes, err = afero.ReadFile(files, "dev.portal.yaml")
	if err == nil {
		return yaml.Unmarshal(bytes, m)
	}

	bytes, err = afero.ReadFile(files, "dev.portal.yml")
	if err == nil {
		return yaml.Unmarshal(bytes, m)
	}

	return
}

func (m *ProjectMetadata) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	return afero.WriteFile(dir, "dev.portal.json", bytes, 0644)
}

func (m *ProjectMetadata) WriteOS(dir string) (err error) {
	return m.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}
