package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/appstore"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"log"
)

type Feat struct {
	port  string
	tray  target.Tray
	spawn target.Spawn
	serve target.Spawn
}

func NewFeat(spawn target.Spawn, tray target.Tray) func(context.Context, bool) error {
	handlers := rpc.Handlers{
		"ping":      func() {},
		"open":      spawn,
		"observe":   appstore.Observe,
		"install":   apps.Install,
		"uninstall": apps.Uninstall,
	}
	return Feat{
		port:  "portal",
		tray:  tray,
		spawn: spawn,
		serve: serve.NewRunner(handlers).Run,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	tray bool,
) (err error) {
	if err = rpc.Command(rpc.NewRequest(id.Anyone, f.port), "ping"); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		if err = f.serve(ctx, f.port); err != nil {
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
