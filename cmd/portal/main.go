package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/clir"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/install"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/uninstall"
	"github.com/cryptopunkscc/go-astral-js/feat/version"
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

	if err := clir.RunPortal(ctx,
		newRuntimeFactory(ctx),
		open.Run,
		apps.List,
		install.Run,
		uninstall.Run,
		version.Run,
	); err != nil {
		cancel()
		log.Fatal(err)
	}

	if ctx.Err() == nil {
		cancel()
		time.Sleep(200 * time.Millisecond)
	}
}

type Adapter struct{ apphost.Flat }

func newRuntimeFactory(ctx context.Context) func(t target.Type, prefix ...string) runtime.Api {
	return func(t target.Type, prefix ...string) runtime.Api {
		switch {
		case t.Is(target.Frontend):
			return &Adapter{Flat: apphost.NewAdapter(ctx, portal.SrvOpenerCtx, prefix...)}
		default:
			return apphost.WithTimeout(ctx, portal.SrvOpenerCtx, prefix...)
		}
	}
}
