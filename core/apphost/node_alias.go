package apphost

import (
	"github.com/cryptopunkscc/astrald/mod/dir/client"
)

func (a *Adapter) NodeAlias() (alias string, err error) {
	client := dir.New(a.TargetID, a.Client)
	identity, err := client.ResolveIdentity(nil, "localnode")
	if err != nil {
		return
	}
	return client.GetAlias(nil, identity)
}
