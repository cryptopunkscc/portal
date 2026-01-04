package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/astrald"
)

func (a *Adapter) NodeAlias() (alias string, err error) {
	dir := astrald.NewDirClient(a.Client)
	identity, err := dir.ResolveIdentity("localnode")
	if err != nil {
		return
	}
	return dir.GetAlias(identity)
}
