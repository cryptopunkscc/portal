package runtime

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/mobile"
	"io"
)

func Apphost(client apphost.Client) mobile.Apphost {
	return &apphostAdapter{client}
}

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
	cc := &apphostConnAdapter{}
	cc.Conn, err = a.i.Query(i, query)
	c = cc
	return
}

var _ io.ReadWriteCloser = &apphostConnAdapter{}

type apphostConnAdapter struct {
	Conn io.ReadWriteCloser
}

func (c *apphostConnAdapter) ReadN(n int) (b []byte, err error) {
	var l int
	b = make([]byte, n)
	if l, err = c.Conn.Read(b); err == nil {
		b = b[:l]
	}
	return
}

func (c *apphostConnAdapter) Read(p []byte) (n int, err error) {
	return c.Conn.Read(p)
}

func (c *apphostConnAdapter) Write(p []byte) (n int, err error) {
	return c.Conn.Write(p)
}

func (c *apphostConnAdapter) Close() error {
	return c.Conn.Close()
}

type apphostListenerAdapter struct{ apphost.Listener }

func (a *apphostListenerAdapter) Next() (query mobile.QueryData, err error) {
	var q apphost.QueryData
	q, err = a.Listener.Next()
	if err != nil {
		return
	}
	query = &queryDataAdapter{query: q}
	return
}

func (a *apphostListenerAdapter) Close() error {
	return a.Listener.Close()
}

type queryDataAdapter struct{ query apphost.QueryData }

func (q *queryDataAdapter) Caller() string {
	return q.query.RemoteIdentity().String()
}

func (q *queryDataAdapter) Accept() (c mobile.Conn, err error) {
	cc := &apphostConnAdapter{}
	cc.Conn, err = q.query.Accept()
	c = cc
	return
}

func (q *queryDataAdapter) Reject() error {
	return q.query.Reject()
}

func (q *queryDataAdapter) Query() string {
	return q.query.Query()
}

func init() {
	astral.Client = *astral.NewClient("memu:apphost", "")
}
