package bind

import (
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/bind"
)

type Apphost interface {
	bind.Apphost
	apphost.Cache
}
