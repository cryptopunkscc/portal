package apphost

import (
	"bufio"
	"github.com/cryptopunkscc/astrald/astral"
	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/google/uuid"
)

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
