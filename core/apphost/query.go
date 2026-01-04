package apphost

import (
	"github.com/cryptopunkscc/astrald/astral"
	lib "github.com/cryptopunkscc/astrald/lib/astrald"
	api "github.com/cryptopunkscc/portal/api/apphost"
)

type query struct{ i *lib.PendingQuery }

func (q *query) Query() string                    { return q.i.Query() }
func (q *query) RemoteIdentity() *astral.Identity { return q.i.Caller() }
func (q *query) Reject() error                    { return q.i.Reject() }
func (q *query) Accept() (api.Conn, error)        { return inConn(q.i.Accept(), nil) }
func (q *query) Close() error                     { return q.i.Close() }
