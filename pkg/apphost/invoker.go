package apphost

import (
	"context"
	"errors"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"log"
	"strings"
	"time"
)

type Invoker struct {
	Flat
	ctx       context.Context
	cancel    context.CancelFunc
	processes sig.Map[string, any]
	invoke    Invoke
}

func NewInvoker(
	ctx context.Context,
	flat Flat,
	serve Invoke,
) (i *Invoker) {
	i = &Invoker{Flat: flat, invoke: serve}
	i.ctx, i.cancel = context.WithCancel(ctx)
	return
}

func (inv *Invoker) Close() error {
	inv.cancel()
	return nil
}

func (inv *Invoker) Query(identity string, query string) (data string, err error) {
	data, err = inv.Flat.Query(identity, query)
	if err != nil && identity == "" {
		if inv.invoke != nil {
			if err := inv.invokeApp(query); err != nil && !errors.Is(err, ErrServiceAlreadyRunning) {
				log.Println("Invoker.Query", inv.port(query), "service not available:", err)
				return data, err
			} else if err == nil {
				log.Println("invoked app for:", query)
			}
		}

		data, err = exec.RetryT[string](8188*time.Millisecond, func(i, n int, d time.Duration) (string, error) {
			if i == 0 {
				return data, err
			}
			return inv.Flat.Query(identity, query)
		})
		if err == nil {
			log.Println("query succeed")
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

	run, err := inv.invoke(inv.Prefix()...)(inv.ctx)
	if err != nil {
		return
	}

	go func() {
		run(src)
		inv.processes.Delete(src)
	}()
	return
}

type Invoke func(prefix ...string) func(ctx context.Context) (func(query string), error)

var ErrServiceAlreadyRunning = errors.New("service already running")
