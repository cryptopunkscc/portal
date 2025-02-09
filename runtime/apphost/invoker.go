package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"time"
)

type Invoker struct {
	apphost.Client
	Invoke target.Run[string]
	Log    plog.Logger
	Ctx    context.Context
}

func (i Invoker) Query(target string, method string, args any) (conn apphost.Conn, err error) {
	conn, err = i.Client.Query(target, method, args)
	if err == nil {
		return
	}

	i.Log.Println("invoking app for:", target, method, args)
	if err = i.Invoke(i.Ctx, method); err != nil {
		return
	}

	conn, err = flow.RetryT[apphost.Conn](
		i.Ctx, 8188*time.Millisecond,
		func(ii, n int, d time.Duration) (apphost.Conn, error) {
			i.Log.Printf("retry query: %s - %d/%d attempt %v: retry after %v", conn, ii+1, n, err, d)
			return i.Client.Query(target, method, args)
		},
	)
	if err == nil {
		i.Log.Println("query succeed", conn.Query())
		return
	}
	return
}
