package apphost

import (
	"github.com/cryptopunkscc/astrald/mod/dir/client"
)

func (a *Adapter) NodeAlias() (alias string, err error) {
	return dir.New(a.HostID(), a.Client).GetAlias(nil, a.HostID())
}
