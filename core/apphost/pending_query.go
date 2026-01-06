package apphost

import (
	"bufio"

	"github.com/cryptopunkscc/astrald/lib/astrald"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/google/uuid"
)

type PendingQuery struct{ *astrald.PendingQuery }

var _ apphost.PendingQuery = &PendingQuery{}

func (q *PendingQuery) Accept() (apphost.Conn, error) {
	c := q.PendingQuery.Accept()
	return &Conn{
		Conn:  c,
		buf:   bufio.NewReader(c),
		ref:   uuid.New().String(),
		query: q.Query(),
		in:    true,
	}, nil
}
