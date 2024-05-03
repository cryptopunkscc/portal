package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("starting portal", os.Args)
	ctx, cancel := context.WithCancel(context.Background())
	go exec.OnShutdown(cancel)
	newRuntime := newRuntimeFactory(ctx)
	clir.Run(ctx, newRuntime)
	if ctx.Err() == nil {
		cancel()
		time.Sleep(200 * time.Millisecond)
	}
}

type Adapter struct{ apphost.Flat }

func newRuntimeFactory(ctx context.Context) func(t target.Type) runtime.Api {
	return func(t target.Type) runtime.Api {
		switch t {
		case target.Frontend:
			return &Adapter{Flat: apphost.NewAdapter(ctx, portal.SrvOpenerCtx)}
		default:
			return apphost.WithTimeout(ctx, portal.SrvOpenerCtx)
		}
	}
}
