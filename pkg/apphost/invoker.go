package apphost

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"io"
	"log"
	"strings"
	"time"
)

type Invoker struct {
	Flat
	ctx       context.Context
	processes sig.Map[io.Closer, interface{}]
}

func (inv *Invoker) Close() error {
	for _, closer := range inv.processes.Keys() {
		_ = closer.Close()
	}
	return nil
}

func (inv *Invoker) Query(identity string, query string) (data string, err error) {
	data, err = inv.Flat.Query(identity, query)
	if err != nil && identity == "" {
		if err := inv.invoke(query); err != nil {
			log.Println("service not available:", err)
			return data, err
		}
		log.Println("invoked app for:", query)

		data, err = exec.Retry[string](8188*time.Millisecond, func(i, n int, d time.Duration) (string, error) {
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

func (inv *Invoker) invoke(query string) (err error) {
	src := strings.Split(query, "?")[0]
	run, closer, err := portal.Bind(src)
	if err != nil {
		return
	}
	go func() {
		inv.processes.Set(closer, 0)
		defer inv.processes.Delete(closer)
		_ = run()
	}()
	return
}
