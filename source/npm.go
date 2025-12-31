package source

import (
	"encoding/json"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/spf13/afero"
)

type NpmProject struct {
	Project
	PackageJson PackageJson
}

func (p *NpmProject) ReadSource(source Source) (err error) {
	p.Source = source
	return p.ReadFs(source.Fs)
}

func (p *NpmProject) ReadFs(files afero.Fs) (err error) {
	return FSReaders{&p.Project, &p.PackageJson}.ReadFs(files)
}

func (p *NpmProject) WriteFs(dir afero.Fs) (err error) {
	return FsWriters{&p.Project, &p.PackageJson}.WriteFs(dir)
}

func (p *NpmProject) WriteOS(dir string) (err error) {
	return p.WriteFs(afero.NewBasePathFs(afero.NewOsFs(), dir))
}

type PackageJson struct {
	Portal  string `json:"portal,omitempty"`
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func (p *PackageJson) IsPortalLib() bool {
	return p.Portal == "lib"
}

func (p *PackageJson) CanBuild() bool {
	return p.Scripts.Build != ""
}

func (p *PackageJson) ReadFs(files afero.Fs) (err error) {
	defer plog.TraceErr(&err)

	bytes, err := afero.ReadFile(files, "package.json")
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, p); err != nil {
		return
	}

	return
}

func (p *PackageJson) WriteFs(dir afero.Fs) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(p)
	if err != nil {
		return
	}
	if err = afero.WriteFile(dir, "package.json", bytes, 0644); err != nil {
		return
	}
	return
}
