package apphost

import (
	"os"
	"path"
	"reflect"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/api/env"
	os2 "github.com/cryptopunkscc/portal/pkg/os"
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
	return a.Client != nil
}

func (a *Adapter) Reconnect() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.connect()
}

func (a *Adapter) connect() (err error) {
	defer plog.TraceErr(&err)
	a.Endpoint = FirstNotZero(a.Endpoint, env.ApphostAddr.Get(), apphost.DefaultEndpoint)
	if a.Token = FirstNotZero(a.Token, env.ApphostToken.Get()); a.Token == "" {
		a.Token = ResolveTokenFromFile()
	}
	host, err := apphost.Connect(a.Endpoint)
	if err != nil {
		return
	}
	if err = host.AuthToken(a.Token); err != nil {
		return
	}
	defer host.Close()
	a.TargetID = host.HostID()
	a.Client = astrald.New(apphost.NewRouter(a.Endpoint, a.Token))
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

func ResolveTokenFromFile(dir ...string) string {
	abs := os2.Abs(dir...)
	file, err := os.Open(path.Join(abs, "astral_user"))
	if err != nil {
		return ""
	}
	defer file.Close()
	o, _, err := astral.Decode(file, astral.Canonical())
	if err != nil {
		return ""
	}
	return o.(*user.CreatedUserInfo).AccessToken.String()
}
