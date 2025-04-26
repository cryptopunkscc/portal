package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"os"
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
	return a.HostID != nil &&
		a.GuestID != nil &&
		a.AuthToken != "" &&
		a.Endpoint != ""
}

func (a *Adapter) Reconnect() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connect()
}

func (a *Adapter) connect() (err error) {
	defer plog.TraceErr(&err)
	if len(a.Endpoint) == 0 {
		a.Endpoint = env.ApphostAddr.Get()
		if len(a.Endpoint) == 0 {
			a.Endpoint = apphost.DefaultEndpoint
		}
	}
	if len(a.AuthToken) == 0 {
		a.AuthToken = os.Getenv(apphost.AuthTokenEnv)
	}
	client, err := apphost.NewClient(a.Endpoint, a.AuthToken)
	if err == nil {
		a.Lib.Client = *client
	}
	return
}
