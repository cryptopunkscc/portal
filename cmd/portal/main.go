package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/apphost"
	"github.com/cryptopunkscc/go-astral-js/pkg/clir"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go exec.OnShutdown(cancel)
	clir.Run(ctx, func() runtime.Api {
		return &Adapter{Flat: apphost.NewAdapter(ctx)}
	})
	if ctx.Err() == nil {
		cancel()
		time.Sleep(200 * time.Millisecond)
	}
}

type Adapter struct{ apphost.Flat }
