package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"log"
)

type Feat struct {
	spawn runtime.Spawn
	tray  runtime.Tray
}

func NewFeat(spawn runtime.Spawn, tray runtime.Tray) func(context.Context, bool) error {
	return Feat{spawn: spawn, tray: tray}.Run
}

func (f Feat) Run(
	ctx context.Context,
	tray bool,
) (err error) {
	port := "portal"
	if err = rpc.Command(rpc.NewRequest(id.Anyone, port), "ping"); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		handlers := rpc.Handlers{
			"ping":      func() {},
			"open":      f.spawn,
			"observe":   appstore.Observe,
			"install":   apps.Install,
			"uninstall": apps.Uninstall,
		}
		if err = serve.Run(ctx, port, handlers); err != nil {
			log.Printf("serve exit: %v\n", err)
		} else {
			log.Println("serve exit")
		}
	}()
	if tray {
		go func() {
			defer cancel()
			f.tray(ctx)
		}()
	}
	<-ctx.Done()
	return
}
