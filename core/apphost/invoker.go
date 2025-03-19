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
	if conn, err = i.Client.Query(target, method, args); err == nil {
		return
	}

	i.Log.Println("invoking app for:", target, method, args)
	if err = i.Invoke(i.Ctx, target); err != nil {
		return
	}

	await := flow.Await{
		UpTo: 8 * time.Millisecond,
		Mod:  8,
		Ctx:  i.Ctx,
	}
	for in := range await.Chan() {
		i.Log.Println("retry:", in)
		if conn, err = i.Client.Query(target, method, args); err == nil {
			i.Log.Println("query succeed", conn.Query())
			return
		}
	}
	return
}
