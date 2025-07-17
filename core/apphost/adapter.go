package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/sig"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"sync"
)

var Default = &Adapter{}

type Adapter struct {
	Lib
	mu         sync.Mutex
	Log        plog.Logger
	identities sig.Map[string, *astral.Identity]
}

func (a *Adapter) Clone() (c *Adapter) {
	c = &Adapter{}
	c.Log = a.Log
	c.Lib = a.Lib
	return
}

func (a *Adapter) Protocol() string {
	if a.Connect() != nil {
		return ""
	}
	return a.Lib.Protocol()
}

func (a *Adapter) Resolve(name string) (i *astral.Identity, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	if name == "" {
		return a.Lib.HostID, nil
	}
	i, ok := a.identities.Get(name)
	if ok {
		return
	}
	if i, err = a.Lib.LocalNode().ResolveIdentity(name); err != nil {
		return
	}
	a.identities.Set(name, i)
	return
}

func (a *Adapter) DisplayName(identity *astral.Identity) string {
	if err := a.Connect(); err != nil {
		return ""
	}
	alias, _ := a.Lib.LocalNode().GetAlias(identity)
	return alias
}

func (a *Adapter) Query(target string, method string, args any) (conn api.Conn, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	id, err := a.Resolve(target)
	if err != nil {
		return
	}
	return outConn(a.Lib.Query(id.String(), method, args))
}

func (a *Adapter) Session() (s api.Session, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	ss, err := a.Lib.Session()
	if err != nil {
		return
	}
	return session{ss}, nil
}

func (a *Adapter) Register() (l api.Listener, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	ll, err := a.Lib.Listen()
	if err != nil {
		return
	}
	l = &listener{i: ll}
	return
}
