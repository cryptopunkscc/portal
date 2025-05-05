package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"io/fs"
)

type Config struct {
	Timeout int64 `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Hidden  bool  `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}

func (c *Config) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, c) }
func (c *Config) LoadFrom(fs fs.FS) error          { return all.Unmarshalers.Load(c, fs, AppFilename) }
