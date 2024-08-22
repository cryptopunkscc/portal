package runtime

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/mobile"
	"io"
)

func Apphost(client apphost.Client) mobile.Apphost { return &apphostAdapter{client} }

type apphostAdapter struct{ i apphost.Client }

func (a *apphostAdapter) Resolve(name string) (s string, err error) {
	var i id.Identity
	i, err = a.i.Resolve(name)
	if err != nil {
		return
	}
	s = i.String()
	return
}
func (a *apphostAdapter) Register(name string) (mobile.ApphostListener, error) {
	listener, err := a.i.Register(name)
	if err != nil {
		return nil, err
	}
	return &apphostListenerAdapter{listener}, err
}
func (a *apphostAdapter) Query(nodeID string, query string) (c mobile.Conn, err error) {
	i, err := id.ParsePublicKeyHex(nodeID)
	if err != nil {
		return
	}
	conn, err := a.i.Query(i, query)
	c = &apphostConnAdapter{conn, reader{conn}}
	return
}

type apphostConnAdapter struct {
	io.WriteCloser
	mobile.Reader
}

type apphostListenerAdapter struct{ apphost.Listener }

func (a *apphostListenerAdapter) Close() error { return a.Listener.Close() }
func (a *apphostListenerAdapter) Next() (query mobile.QueryData, err error) {
	var q apphost.QueryData
	q, err = a.Listener.Next()
	if err != nil {
		return
	}
	query = &queryDataAdapter{query: q}
	return
}

type queryDataAdapter struct{ query apphost.QueryData }

func (q *queryDataAdapter) Query() string  { return q.query.Query() }
func (q *queryDataAdapter) Caller() string { return q.query.RemoteIdentity().String() }
func (q *queryDataAdapter) Reject() error  { return q.query.Reject() }
func (q *queryDataAdapter) Accept() (c mobile.Conn, err error) {
	conn, err := q.query.Accept()
	c = &apphostConnAdapter{conn, reader{conn}}
	return
}

func init() {
	astral.Client = *astral.NewClient("memu:apphost", "")
}
