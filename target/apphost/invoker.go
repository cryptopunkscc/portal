package apphost

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/target"
	"strings"
	"time"
)

type Invoker struct {
	*Adapter
	ctx       context.Context
	cancel    context.CancelFunc
	processes sig.Map[string, any]
	invoke    target.Dispatch
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
	data, err = inv.Adapter.Query(identity, query)
	if err != nil && identity == "" {
		if inv.invoke != nil {
			if err := inv.invokeApp(query); err != nil && !errors.Is(err, ErrServiceAlreadyRunning) {
				err = fmt.Errorf("Invoker.Query %s service not available: %v", query, err)
				return data, err
			} else if err == nil {
				inv.log.Println("invoked app for:", query)
			}
		}

		data, err = exec.RetryT[string](inv.ctx, 8188*time.Millisecond, func(i, n int, d time.Duration) (string, error) {
			if i == 0 {
				return data, err
			}
			return inv.Adapter.Query(identity, query)
		})
		if err == nil {
			inv.log.Println("query succeed")
			return
		}
	}
	return
}

func (inv *Invoker) invokeApp(query string) (err error) {
	src := strings.Split(query, "?")[0]
	if _, ok := inv.processes.Set(src, 0); !ok {
		return ErrServiceAlreadyRunning
	}

	src = strings.Join(append(inv.Prefix(), src), ".")

	go func() {
		if err := inv.invoke(inv.ctx, src); err != nil {
			inv.log.Println("Invoker.invokeApp:", err)
		}
		inv.processes.Delete(src)
	}()
	return
}

var ErrServiceAlreadyRunning = errors.New("service already running")
