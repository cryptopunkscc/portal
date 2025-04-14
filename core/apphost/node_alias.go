package apphost

import "github.com/cryptopunkscc/portal/pkg/rpc"

func (a *Adapter) NodeAlias() (alias string, err error) {
	p, err := rpc.Query[profile](a.Rpc().Request("localnode"), ".profile")
	if err == nil {
		alias = p.Alias
	}
	return
}

type profile struct {
	Alias string `json:"alias"`
}
