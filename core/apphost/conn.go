package apphost

import (
	"bufio"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
)

type Conn struct {
	astral.Conn
	buf   *bufio.Reader
	query string
	ref   string
	in    bool
}

var _ apphost.Conn = &Conn{}

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
