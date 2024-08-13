package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"strings"
	"time"
)

func Invoker(
	ctx context.Context,
	flat target.Apphost,
	invoke target.Request,
) target.Apphost {
	i := &invoker{Apphost: flat, invoke: invoke}
	i.ctx, i.cancel = context.WithCancel(ctx)
	i.log = plog.Get(i.ctx).Type(i)
	return i
}

type invoker struct {
	target.Apphost
	ctx    context.Context
	cancel context.CancelFunc
	invoke target.Request
	log    plog.Logger
}

func (i *invoker) Close() error {
	i.cancel()
	return nil
}

func (i *invoker) Query(identity string, query string) (data string, err error) {
	data, err = i.Apphost.Query(identity, query)
	if err != nil && identity == "" {
		if i.invoke != nil {
			i.log.Println("invoking app for:", query)
			if err = i.invokeApp(query); err != nil {
				return
			}
		}

		data, err = flow.RetryT[string](i.ctx, 8188*time.Millisecond, func(ii, n int, d time.Duration) (string, error) {
			i.log.Printf("retry query: %s - %d/%d attempt %v: retry after %v", data, ii+1, n, err, d)
			return i.Apphost.Query(identity, query)
		})
		if err == nil {
			i.log.Println("query succeed", data)
			return
		}
	}
	return
}

func (i *invoker) invokeApp(query string) error {
	src := strings.Split(query, "?")[0]
	return i.invoke(i.ctx, src)
}
