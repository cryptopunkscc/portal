package apphost

import (
	"net"

	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Listener struct{ *astrald.Listener }

var _ apphost.Listener = &Listener{}

func (l *Listener) Next() (q apphost.PendingQuery, err error) {
	qq := PendingQuery{}
	if qq.PendingQuery, err = l.Listener.Next(); err != nil {
		return
	}
	return &qq, nil
}

func (l *Listener) Accept() (net.Conn, error) { return l.Listener.Accept() }
func (l *Listener) Close() (err error) {
	plog.TraceErr(&err)
	return l.Listener.Close()
}
