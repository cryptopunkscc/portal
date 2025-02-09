package apphost

import (
	"bufio"
	"context"
	"github.com/cryptopunkscc/astrald/astral"
	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	mod "github.com/cryptopunkscc/astrald/mod/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/google/uuid"
	"net"
	"sync"
)

var Default = Adapter(Lib)

func Adapter(client *lib.Client) api.Client { return &adapter{lib: client} }

type adapter struct {
	mu  sync.RWMutex
	lib *lib.Client
}

func (a *adapter) client() *lib.Client {
	if a.lib == nil {
		a.lib = Lib
	}
	if a.lib == Lib && !IsConnected() {
		a.mu.Lock()
		defer a.mu.Unlock()
		if IsConnected() {
			return a.lib
		}
		ctx := context.Background()
		log := plog.New().W().Type(a).Set(&ctx)
		log.Println("apphost not connected")
		if err := Connect(ctx); err != nil {
			log.P().Println(err)
		}
	}
	return a.lib
}

func (a *adapter) Protocol() string { return a.client().Protocol() }
func (a *adapter) Resolve(name string) (*astral.Identity, error) {
	if name == "" {
		return a.lib.HostID, nil
	}
	return a.client().ResolveIdentity(name)
}
func (a *adapter) DisplayName(identity *astral.Identity) string {
	return a.client().DisplayName(identity)
}

func (a *adapter) Query(target string, method string, args any) (conn api.Conn, err error) {
	id, err := a.Resolve(target)
	if err != nil {
		return
	}
	return outConn(a.client().Query(id.String(), method, args))
}

func (a *adapter) Session() (api.Session, error) {
	s, err := a.client().Session()
	if err != nil {
		return nil, err
	}
	return session{s}, nil
}

func (a *adapter) Register() (l api.Listener, err error) {
	ll, err := a.client().Listen()
	if err != nil {
		return
	}
	l = &listener{i: ll}
	return
}

type session struct{ i *lib.Session }

func (s session) Token(token string) (res api.TokenResponse, err error) {
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
func (l *listener) Close() (err error)        { return l.i.Close() }
func (l *listener) Addr() net.Addr            { return l.i.Addr() }

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
