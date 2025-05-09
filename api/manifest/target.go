package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/dec/all"
	"runtime"
)

type Target struct {
	Exec string `json:"exec,omitempty" yaml:"exec,omitempty"`
	OS   string `json:"os,omitempty" yaml:"os,omitempty"`
	Arch string `json:"arch,omitempty" yaml:"arch,omitempty"`
}

func (r *Target) UnmarshalFrom(bytes []byte) error { return all.Unmarshalers.Unmarshal(bytes, r) }

func (r Target) Match() bool { return r.OS == runtime.GOOS && r.Arch == runtime.GOARCH }
