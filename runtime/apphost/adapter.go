package apphost

import (
	"bufio"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/astrald/lib/astral"
	"github.com/cryptopunkscc/astrald/mod/apphost/proto"
	. "github.com/cryptopunkscc/portal/api/apphost"
	"net"
)

func Default() Cached {
	return Adapter(astral.Client)
}

func Adapter(client astral.ApphostClient) Cached {
	return &adapter{
		i:     client,
		cache: newCache(),
	}
}

type adapter struct {
	i astral.ApphostClient
	*cache
}

func (a adapter) Interrupt() {
	for _, closer := range a.listeners.Release() {
		_ = closer.Close()
	}
	for _, closer := range a.connections.Release() {
		_ = closer.Close()
	}
}

func (a adapter) Session() (s Session, err error) {
	ss := &session{}
	if ss.i, err = a.i.Session(); err != nil {
		return
	}
	ss.cache = a.cache
	return ss, nil
}
func (a adapter) Discovery() Discovery {
	return a.i.Discovery()
}
func (a adapter) Query(remoteID id.Identity, query string) (conn Conn, err error) {
	return a.outConn(a.i.Query(remoteID, query))
}
func (a adapter) QueryName(name string, query string) (conn Conn, err error) {
	return a.outConn(a.i.QueryName(name, query))
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
func (a adapter) Register(service string) (l Listener, err error) {
	ll, err := a.i.Register(service)
	return a.setListener(service, ll, err)
}

type session struct {
	i *astral.Session
	*cache
}

func (s session) Query(remoteID id.Identity, query string) (Conn, error) {
	return s.outConn(s.i.Query(remoteID, query))
}
func (s session) Resolve(name string) (id.Identity, error) {
	return s.i.Resolve(name)
}
func (s session) NodeInfo(identity id.Identity) (proto.NodeInfoData, error) {
	return s.i.NodeInfo(identity)
}
func (s session) Register(service string, target string) error { return s.i.Register(service, target) }
func (s session) Exec(identity id.Identity, app string, args []string, env []string) (err error) {
	return s.i.Exec(identity, app, args, env)
}

type listener struct {
	i *astral.Listener
	*cache
	port string
}

func (l *listener) Port() string { return l.port }
func (l *listener) Next() (q QueryData, err error) {
	qq := query{cache: l.cache}
	if qq.i, err = l.i.Next(); err != nil {
		return
	}
	return &qq, nil
}
func (l *listener) QueryCh() (c <-chan QueryData) {
	in := l.i.QueryCh()
	out := make(chan QueryData)
	go func() {
		for data := range in {
			out <- &query{data, l.cache}
		}
	}()
	return out
}
func (l *listener) Accept() (net.Conn, error)  { return l.i.Accept() }
func (l *listener) AcceptAll() <-chan net.Conn { return l.i.AcceptAll() }
func (l *listener) Close() (err error) {
	err = l.i.Close()
	l.cache.deleteListener(l.port)
	return
}
func (l *listener) Addr() net.Addr { return l.i.Addr() }
func (l *listener) Target() string { return l.i.Target() }

type query struct {
	i *astral.QueryData
	*cache
}

func (q *query) Query() string               { return q.i.Query() }
func (q *query) RemoteIdentity() id.Identity { return q.i.RemoteIdentity() }
func (q *query) Reject() error               { return q.i.Reject() }
func (q *query) Accept() (Conn, error)       { return q.inConn(q.i.Accept()) }

type conn struct {
	*astral.Conn
	*cache
	buf *bufio.Reader
	ref string
	in  bool
}

func (c *conn) Read(b []byte) (n int, err error)      { return c.buf.Read(b) }
func (c *conn) ReadString(delim byte) (string, error) { return c.buf.ReadString(delim) }
func (c *conn) Ref() string                           { return c.ref }
func (c *conn) In() bool                              { return c.in }
func (c *conn) Close() error {
	c.cache.deleteConn(c)
	return c.Conn.Close()
}
