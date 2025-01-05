package start

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/portal"
)

func Create(deps Deps) Start { return New(deps.Portal(), deps.Apphost()) }

type Deps interface {
	Portal() portal.Client
	Apphost() apphost.Client
}
