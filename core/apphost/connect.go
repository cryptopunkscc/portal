package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func (a *Adapter) Connect() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.IsConnected() {
		return
	}
	return a.connect()
}

func (a *Adapter) IsConnected() bool {
	return a.HostID != nil
}

func (a *Adapter) Reconnect() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connect()
}

func (a *Adapter) connect() (err error) {
	defer plog.TraceErr(&err)
	defaultConfig := astrald.DefaultConfig()
	if len(a.Endpoint) == 0 {
		a.Endpoint = defaultConfig.Endpoint
	}
	if len(a.Token) == 0 {
		a.Token = defaultConfig.Token
	}
	a.Client = astrald.NewClient(a.Config)
	if a.HostID == nil {
		a.HostID, err = astrald.NewDirClient(a.Client).ResolveIdentity("localnode")
	}
	return
}
