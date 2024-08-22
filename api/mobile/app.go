package mobile

import (
	"github.com/cryptopunkscc/portal/api/bind"
)

type App interface {
	Manifest() *Manifest
	Assets() Assets
	Runtime() bind.Runtime
}
