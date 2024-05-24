package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Feat struct {
	port  string
	tray  target.Tray
	serve target.Dispatch
}

func NewFeat(spawn target.Dispatch, tray target.Tray) func(context.Context, bool) error {
	handlers := rpc.Handlers{
		"ping":      func() {},
		"open":      spawn,
		"observe":   apps.Observe,
		"install":   apps.Install,
		"uninstall": apps.Uninstall,
	}
	return Feat{
		port:  "portal",
		tray:  tray,
		serve: serve.NewRunner(handlers).Run,
	}.Run
}

func (f Feat) Run(
	ctx context.Context,
	tray bool,
) (err error) {
	log := plog.Get(ctx).Type(f).Set(&ctx)
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
