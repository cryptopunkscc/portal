package apphost

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/target"
	"strings"
	"time"
)

type Invoker struct {
	*Adapter
	ctx    context.Context
	cancel context.CancelFunc
	invoke target.Dispatch
}

func NewInvoker(
	ctx context.Context,
	flat *Adapter,
	invoke target.Dispatch,
) (i *Invoker) {
	i = &Invoker{Adapter: flat, invoke: invoke}
	i.ctx, i.cancel = context.WithCancel(ctx)
	return
}

func (inv *Invoker) Close() error {
	inv.cancel()
	return nil
}

func (inv *Invoker) Query(identity string, query string) (data string, err error) {
	log := inv.log.Type(inv)
	data, err = inv.Adapter.Query(identity, query)
	if err != nil && identity == "" {
		if inv.invoke != nil {
			log.Println("invoking app for:", query)
			if err = inv.invokeApp(query); err != nil {
				return
			}
		}

		data, err = flow.RetryT[string](inv.ctx, 8188*time.Millisecond, func(i, n int, d time.Duration) (string, error) {
			log.Printf("retry query: %s - %d/%d attempt %v: retry after %v", data, i+1, n, err, d)
			return inv.Adapter.Query(identity, query)
		})
		if err == nil {
			log.Println("query succeed", data)
			return
		}
	}
	return
}

func (inv *Invoker) invokeApp(query string) error {
	src := strings.Split(query, "?")[0]
	return inv.invoke(inv.ctx, src)
}
