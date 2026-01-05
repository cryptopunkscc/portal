package npm

import (
	"encoding/json"
	"path"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/source"
	"github.com/spf13/afero"
)

type PackageJson struct {
	Scripts struct {
		Build string `json:"build"`
	} `json:"scripts,omitempty"`
}

func (p *PackageJson) CanBuild() bool {
	return p.Scripts.Build != ""
}

func (p *PackageJson) ReadSrc(src source.Source) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := afero.ReadFile(src.Ref_().Fs, path.Join(src.Ref_().Path, "package.json"))
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, p)
}

func (p *PackageJson) WriteRef(ref source.Ref) (err error) {
	defer plog.TraceErr(&err)
	bytes, err := json.Marshal(p)
	if err != nil {
		return
	}
	if err = afero.WriteFile(ref.Fs, path.Join(ref.Path, "package.json"), bytes, 0644); err != nil {
		return
	}
	return
}
