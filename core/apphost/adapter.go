package apphost

import (
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

var Default = &Adapter{}

type Adapter struct {
	*astrald.Client
	Endpoint string
	Token    string
	TargetID *astral.Identity
	mu       sync.Mutex
	Log      plog.Logger
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
	return a.Dir().ResolveIdentity(nil, name)
}

func (a *Adapter) DisplayName(identity *astral.Identity) string {
	if err := a.Connect(); err != nil {
		return ""
	}
	alias, _ := a.Dir().GetAlias(nil, identity)
	return alias
}

func (a *Adapter) Query(target string, method string, args any) (conn apphost.Conn, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	if err != nil {
		return
	}
	id, err := a.Resolve(target)
	if err != nil {
		return
	}
	return outConn(a.Client.WithTarget(id).Query(nil, method, args))
}

func (a *Adapter) Register() (out apphost.Listener, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	l, err := astrald.NewAppHostClient(a.TargetID, a.Client).RegisterHandler(nil)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	return &Listener{l}, nil
}
