package apphost

import (
	"bufio"
	"github.com/cryptopunkscc/astrald/astral"
	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	"github.com/cryptopunkscc/astrald/sig"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/google/uuid"
	"net"
)

var Default = &Adapter{}

type Adapter struct {
	Client
	Log        plog.Logger
	identities sig.Map[string, *astral.Identity]
}

func (a *Adapter) Connect() (err error) {
	defer plog.TraceErr(&err)
	if !a.Client.IsConnected() {
		return a.Client.Connect()
	}
	return
}

func (a *Adapter) Protocol() string {
	if a.Connect() != nil {
		return ""
	}
	return a.Client.Protocol()
}

func (a *Adapter) Resolve(name string) (i *astral.Identity, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	if name == "" {
		return a.Client.HostID, nil
	}
	i, ok := a.identities.Get(name)
	if ok {
		return
	}
	if i, err = a.Client.ResolveIdentity(name); err != nil {
		return
	}
	a.identities.Set(name, i)
	return
}

func (a *Adapter) DisplayName(identity *astral.Identity) string {
	if err := a.Connect(); err != nil {
		return ""
	}
	return a.Client.DisplayName(identity)
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
	return outConn(a.Client.Query(id.String(), method, args))
}

func (a *Adapter) Session() (s api.Session, err error) {
	defer plog.TraceErr(&err)
	if err = a.Connect(); err != nil {
		return
	}
	ss, err := a.Client.Session()
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
	ll, err := a.Client.Listen()
	if err != nil {
		return
	}
	l = &listener{i: ll}
	return
}

type session struct{ i *lib.Session }

func (s session) Token(token string) (res api.TokenResponse, err error) {
	defer plog.TraceErr(&err)
	response, err := s.i.Token(token)
	if err != nil {
		return nil, err
	}
	return &tokenResponse{response}, nil
}

func (s session) Query(callerID *astral.Identity, targetID *astral.Identity, query string) (api.Conn, error) {
	return outConn(s.i.Query(callerID, targetID, query))
}

func (s session) Register(identity *astral.Identity, target string) (token string, err error) {
	return s.i.Register(identity, target)
}

func (s session) Close() error {
	return s.i.Close()
}

type tokenResponse struct{ i mod.TokenResponse }

func (t tokenResponse) Code() uint8               { return uint8(t.i.Code) }
func (t tokenResponse) GuestID() *astral.Identity { return t.i.GuestID }
func (t tokenResponse) HostID() *astral.Identity  { return t.i.HostID }

type listener struct{ i *lib.Listener }

func (l *listener) String() string        { return l.i.String() }
func (l *listener) Token() string         { return l.i.Token() }
func (l *listener) SetToken(token string) { l.i.SetToken(token) }
func (l *listener) Done() <-chan struct{} { return l.i.Done() }

func (l *listener) Next() (q api.PendingQuery, err error) {
	qq := query{}
	if qq.i, err = l.i.Next(); err != nil {
		return
	}
	return &qq, nil
}

func (l *listener) Accept() (net.Conn, error) { return l.i.Accept() }
func (l *listener) Close() (err error) {
	plog.TraceErr(&err)
	return l.i.Close()
}
func (l *listener) Addr() net.Addr { return l.i.Addr() }

type query struct{ i *lib.PendingQuery }

func (q *query) Query() string                    { return q.i.Query() }
func (q *query) RemoteIdentity() *astral.Identity { return q.i.Caller() }
func (q *query) Reject() error                    { return q.i.Reject() }
func (q *query) Accept() (api.Conn, error)        { return inConn(q.i.Accept()) }
func (q *query) Close() error                     { return q.i.Close() }

type conn struct {
	*lib.Conn
	buf *bufio.Reader
	ref string
	in  bool
}

var _ api.Conn = &conn{}

func inConn(c *lib.Conn, err error) (*conn, error)  { return newConn(c, err, true) }
func outConn(c *lib.Conn, err error) (*conn, error) { return newConn(c, err, false) }
func newConn(c *lib.Conn, err error, in bool) (*conn, error) {
	defer plog.TraceErr(&err)
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

func (c *conn) RemoteIdentity() *astral.Identity {
	return c.Conn.RemoteIdentity()
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
