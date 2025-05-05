package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"io/fs"
)

type Release struct {
	Version int `json:"version,omitempty" yaml:"version,omitempty"`
}

func (r *Release) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, r) }
func (r *Release) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(r, fs, AppFilename) }
