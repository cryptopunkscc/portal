package manifest

import (
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

type Api struct {
	Version     int `json:"version,omitempty" yaml:"version,omitempty"`
	cmd.Handler `json:",omitempty" yaml:",omitempty,inline"`
}
