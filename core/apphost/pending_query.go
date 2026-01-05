package apphost

import (
	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/api/apphost"
)

type PendingQuery struct{ *astrald.PendingQuery }

var _ apphost.PendingQuery = &PendingQuery{}

func (q *PendingQuery) Accept() (apphost.Conn, error) { return inConn(q.PendingQuery.Accept(), nil) }
