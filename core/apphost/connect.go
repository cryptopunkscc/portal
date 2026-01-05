package apphost

import (
	"reflect"

	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/api/env"
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
	return a.TargetID != nil
}

func (a *Adapter) Reconnect() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connect()
}

func (a *Adapter) connect() (err error) {
	defer plog.TraceErr(&err)
	a.Endpoint = FirstNotZero(a.Endpoint, env.ApphostAddr.Get(), apphost.DefaultEndpoint)
	a.Token = FirstNotZero(a.Token, env.ApphostToken.Get())
	host, err := apphost.Connect(a.Endpoint)
	if err != nil {
		return
	}
	if err = host.AuthToken(a.Token); err != nil {
		return
	}
	defer host.Close()
	a.TargetID = host.HostID()
	a.Client = astrald.NewClient(apphost.NewRouter(a.Endpoint, a.Token))
	return
}

func FirstNotZero[T any](anyOf ...T) (zero T) {
	for _, next := range anyOf {
		if val := reflect.ValueOf(next); !val.IsZero() {
			return next
		}
	}
	return
}
