package apphost

import (
	"bufio"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/google/uuid"
)

type Conn struct {
	astral.Conn
	buf   *bufio.Reader
	query string
	ref   string
	in    bool
}

var _ apphost.Conn = &Conn{}

func inConn(c astral.Conn, err error) (*Conn, error)  { return newConn(c, err, true) }
func outConn(c astral.Conn, err error) (*Conn, error) { return newConn(c, err, false) }

func newConn(c astral.Conn, err error, in bool) (*Conn, error) {
	defer plog.TraceErr(&err)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Conn: c,
		buf:  bufio.NewReader(c),
		ref:  uuid.New().String(),
		in:   in,
	}, nil
}
func (c *Conn) Query() string { return c.query }

func (c *Conn) Write(b []byte) (n int, err error) {
	if n, err = c.Conn.Write(b); err != nil {
		_ = c.Close()
	}
	return
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if n, err = c.buf.Read(b); err != nil {
		_ = c.Close()
	}
	return
}
func (c *Conn) ReadString(delim byte) (s string, err error) {
	if s, err = c.buf.ReadString(delim); err != nil {
		_ = c.Close()
	}
	return
}
func (c *Conn) Ref() string { return c.ref }
func (c *Conn) In() bool    { return c.in }
