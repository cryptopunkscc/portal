package apphost

import (
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

var Default = &Adapter{}

type Adapter struct {
	*astrald.Client
	astrald.Config
	HostID *astral.Identity
	mu     sync.Mutex
	Log    plog.Logger
}

func (a *Adapter) Clone() (c *Adapter) {
	c = &Adapter{}
	c.Log = a.Log
	return
}

func (a *Adapter) Resolve(name string) (i *astral.Identity, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	if name == "" {
		name = "localnode"
	}
	return astrald.NewDirClient(a.Client).ResolveIdentity(name)
}

func (a *Adapter) DisplayName(identity *astral.Identity) string {
	if err := a.Connect(); err != nil {
		return ""
	}
	alias, _ := astrald.NewDirClient(a.Client).GetAlias(identity)
	return alias
}

func (a *Adapter) Query(target string, method string, args any) (conn api.Conn, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	if err != nil {
		return
	}
	return outConn(a.Client.Query(target, method, args))
}

func (a *Adapter) Register() (l api.Listener, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	ll, err := a.Client.Listen()
	if err != nil {
		return
	}
	l = &listener{i: ll}
	return
}
