package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"io/fs"
)

type Dist struct {
	App     `json:",inline" yaml:",inline"`
	Api     Api     `json:"api,omitempty" yaml:"api,omitempty"`
	Config  Config  `json:"config,omitempty" yaml:"config,omitempty"`
	Target  Target  `json:"target,omitempty" yaml:"target,omitempty"`
	Release Release `json:"release,omitempty" yaml:"release,omitempty"`
}

func (d *Dist) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, d) }
func (d *Dist) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(d, fs, AppFilename) }
