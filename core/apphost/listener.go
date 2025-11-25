package apphost

import (
	"bufio"
	"net"

	lib "github.com/cryptopunkscc/astrald/lib/apphost"
	api "github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

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

type conn struct {
	*lib.Conn
	buf *bufio.Reader
	ref string
	in  bool
}
