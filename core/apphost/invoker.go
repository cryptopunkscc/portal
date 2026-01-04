package apphost

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/rpc"
)

type Invoker struct {
	Cached
	Ctx context.Context
}

func (i *Invoker) Query(target string, method string, args any) (conn apphost.Conn, err error) {
	if conn, err = i.Cached.Query(target, method, args); err == nil {
		return
	}

	i.Log.Println("invoking app for:", target, method, args)
	if err = i.Open(target); err != nil {
		return
	}

	await := flow.Await{
		UpTo: 8 * time.Second,
		Mod:  8,
		Ctx:  i.Ctx,
	}
	for in := range await.Chan() {
		i.Log.Println("retry:", in)
		if conn, err = i.Cached.Query(target, method, args); err == nil {
			i.Log.Println("query succeed", conn.Query())
			return
		}
	}
	return
}

func (i *Invoker) Open(src string) (err error) {
	i.Log.Println("starting query", "portald.open", src)
	request := i.Cached.Rpc().Request("portald")
	err = rpc.Command(request, "open", src)
	if err != nil {
		i.Log.E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	i.Log.Println("started query", "portald.open", src)
	return
}
