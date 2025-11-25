package manifest

import (
	"io/fs"

	"github.com/cryptopunkscc/portal/pkg/dec/all"
)

const DevFilename = "dev.portal"

type Dev struct {
	Dist   `json:",inline" yaml:",inline"`
	Builds `json:",inline,omitempty" yaml:",inline,omitempty"`
}

func (d *Dev) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, d) }
func (d *Dev) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(d, fs, DevFilename) }
