package apphost

import (
	"bufio"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/mod/apphost/proto"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/google/uuid"
	"net"
)

func Adapter(client astral.ApphostClient) apphost.Client {
	return &adapter{
		i: client,
	}
}

type adapter struct {
	i astral.ApphostClient
}

func (a adapter) Discovery() apphost.Discovery {
	return a.i.Discovery()
}
func (a adapter) Query(remoteID id.Identity, query string) (apphost.Conn, error) {
	return outConn(a.i.Query(remoteID, query))
}
func (a adapter) QueryName(name string, query string) (conn apphost.Conn, err error) {
	return outConn(a.i.QueryName(name, query))
}
func (a adapter) Resolve(name string) (id.Identity, error) {
	return a.i.Resolve(name)
}
func (a adapter) NodeInfo(identity id.Identity) (info proto.NodeInfoData, err error) {
	return a.i.NodeInfo(identity)
}
func (a adapter) Exec(identity id.Identity, app string, args []string, env []string) error {
	return a.i.Exec(identity, app, args, env)
}
func (a adapter) Register(service string) (l apphost.Listener, err error) {
	ll, err := a.i.Register(service)
	if err != nil {
		return
	}
	l = &listener{
		i:    ll,
		port: service,
	}
	return
}

type listener struct {
	i    *astral.Listener
	port string
}

func (l *listener) Port() string { return l.port }
func (l *listener) Next() (q apphost.QueryData, err error) {
	qq := query{}
	if qq.i, err = l.i.Next(); err != nil {
		return
	}
	return &qq, nil
}
func (l *listener) QueryCh() (c <-chan apphost.QueryData) {
	in := l.i.QueryCh()
	out := make(chan apphost.QueryData)
	go func() {
		for data := range in {
			out <- &query{data}
		}
	}()
	return out
}
func (l *listener) Accept() (net.Conn, error)  { return l.i.Accept() }
func (l *listener) AcceptAll() <-chan net.Conn { return l.i.AcceptAll() }
func (l *listener) Close() (err error)         { return l.i.Close() }
func (l *listener) Addr() net.Addr             { return l.i.Addr() }
func (l *listener) Target() string             { return l.i.Target() }

type query struct{ i *astral.QueryData }

func (q *query) Query() string                 { return q.i.Query() }
func (q *query) RemoteIdentity() id.Identity   { return q.i.RemoteIdentity() }
func (q *query) Reject() error                 { return q.i.Reject() }
func (q *query) Accept() (apphost.Conn, error) { return inConn(q.i.Accept()) }

type conn struct {
	*astral.Conn
	buf *bufio.Reader
	ref string
	in  bool
}

func inConn(c *astral.Conn, err error) (*conn, error)  { return newConn(c, err, true) }
func outConn(c *astral.Conn, err error) (*conn, error) { return newConn(c, err, false) }
func newConn(c *astral.Conn, err error, in bool) (*conn, error) {
	if err != nil {
		return nil, err
	}
	return &conn{
		Conn: c,
		buf:  bufio.NewReader(c),
		ref:  uuid.New().String(),
		in:   in,
	}, nil
}

func (c *conn) Write(b []byte) (n int, err error) {
	if n, err = c.Conn.Write(b); err != nil {
		_ = c.Close()
	}
	return
}

func (c *conn) Read(b []byte) (n int, err error) {
	if n, err = c.buf.Read(b); err != nil {
		_ = c.Close()
	}
	return
}
func (c *conn) ReadString(delim byte) (s string, err error) {
	if s, err = c.buf.ReadString(delim); err != nil {
		_ = c.Close()
	}
	return
}
func (c *conn) Ref() string { return c.ref }
func (c *conn) In() bool    { return c.in }
