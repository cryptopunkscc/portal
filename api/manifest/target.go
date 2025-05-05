package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"io/fs"
)

type Target struct {
	Exec string `json:"exec,omitempty" yaml:"exec,omitempty"`
	OS   string `json:"os,omitempty" yaml:"os,omitempty"`
	Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
}

func (r *Target) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, r) }
func (r *Target) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(r, fs, AppFilename) }
