package apphost

import (
	"bufio"
	"context"
	"sync"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/astrald/lib/query"
	apphost2 "github.com/cryptopunkscc/astrald/mod/apphost/client"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/google/uuid"
)

var Default = &Adapter{}

type Adapter struct {
	adapter
	mu sync.Mutex
}

type adapter struct {
	*astrald.Client
	TargetID *astral.Identity
	Endpoint string
	Token    string
	Log      plog.Logger
}

func (a *Adapter) Clone() (c *Adapter) {
	return &Adapter{adapter: a.adapter}
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
	q := query.New(a.GuestID(), id, method, args)
	aConn, err := a.Client.RouteQuery(astral.NewContext(nil), q)
	if err != nil {
		return
	}
	return &Conn{
		Conn:  aConn,
		buf:   bufio.NewReader(aConn),
		ref:   uuid.New().String(),
		query: q.Query,
	}, nil
}

func (a *Adapter) Register(ctx context.Context) (out apphost.Listener, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}

	client := apphost2.New(a.TargetID, a.Client)

	listener, err := astrald.Listen()
	if err != nil {
		return nil, err
	}

	err = client.RegisterHandler(astral.NewContext(ctx), listener.Endpoint(), listener.AuthToken())
	if err != nil {
		return nil, err
	}

	return &Listener{listener}, nil
}
