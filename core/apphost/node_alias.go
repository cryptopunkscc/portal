package apphost

import "github.com/cryptopunkscc/portal/api/dir"

func (a *Adapter) NodeAlias() (alias string, err error) {
	return dir.Op(a, "localnode").GetAlias(*a.HostID)
}
