package astrald

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
)

type Config struct {
	Node    core.Config    `yaml:",omitempty"`
	Apphost apphost.Config `yaml:",omitempty"`
	Ether   ether.Config   `yaml:",omitempty"`
}
