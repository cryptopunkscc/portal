package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"strings"
	"time"
)

type Invoker struct {
	apphost.Client
	Invoke target.Request
	Log    plog.Logger
	Ctx    context.Context
}

func (i Invoker) Query(identity id.Identity, query string) (conn apphost.Conn, err error) {
	conn, err = i.Client.Query(identity, query)
	if err == nil {
		return
	}
	if identity != id.Anyone {
		return
	}

	i.Log.Println("invoking app for:", query)
	src := strings.Split(query, "?")[0]
	src = strings.TrimPrefix(src, apphost.NewPort("").String()) // hacky way to handle (remove) dev prefix FIXME
	if err = i.Invoke(i.Ctx, src); err != nil {
		return
	}

	conn, err = flow.RetryT[apphost.Conn](
		i.Ctx, 8188*time.Millisecond,
		func(ii, n int, d time.Duration) (apphost.Conn, error) {
			i.Log.Printf("retry query: %s - %d/%d attempt %v: retry after %v", conn, ii+1, n, err, d)
			return i.Client.Query(identity, query)
		},
	)
	if err == nil {
		i.Log.Println("query succeed", conn.Query())
		return
	}
	return
}
