package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Feat struct {
	port  string
	tray  target.Tray
	serve target.Dispatch
}

type (
	Observe   func(ctx context.Context, conn rpc.Conn) (err error)
	Install   func(src string) error
	Uninstall func(id string) error
	Service   func(handlers rpc.Handlers) target.Dispatch
)

func NewFeat(
	port string,
	service Service,
	handlers rpc.Handlers,
	spawn target.Dispatch,
	observe Observe,
	tray target.Tray,
) *Feat {
	if handlers == nil {
		handlers = rpc.Handlers{}
	}
	handlers["ping"] = func() {}
	handlers["open"] = spawn
	handlers["observe"] = observe
	return &Feat{
		port:  port,
		tray:  tray,
		serve: service(handlers),
	}
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

func (f Feat) Dispatch(ctx context.Context, _ string, _ ...string) (err error) {
	go func() {
		if err = f.Run(ctx, false); err != nil {
			plog.Get(ctx).Type(f).Println("dispatch:", err)
		}
	}()
	return
}
