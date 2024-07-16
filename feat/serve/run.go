package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
)

type Deps[T target.Portal] interface {
	Astral() Astral
	Port() target.Port
	RunService() Service
	RpcHandlers() rpc.Handlers
	FeatObserve() Observe
	NewTray() target.NewTray
	RunSpawn() target.Dispatch
}

func Inject[T target.Portal](deps Deps[T]) *Feat {
	var runTray target.Tray
	if deps.NewTray() != nil {
		runTray = deps.NewTray()(deps.RunSpawn())
	}
	return NewFeat(
		deps.Astral(),
		deps.Port(),
		deps.RunService(),
		deps.RpcHandlers(),
		deps.RunSpawn(),
		deps.FeatObserve(),
		runTray,
	)
}

// Feat representing portal service.
type Feat struct {
	astral Astral
	port   target.Port
	tray   target.Tray
	serve  target.Dispatch
}

type (
	// Astral starts daemon if not already running.
	Astral func(ctx context.Context) (err error)

	// Observe on installed applications.
	Observe func(ctx context.Context, conn rpc.Conn) (err error)

	// Service starts portal rpc with additional handlers.
	Service func(handlers rpc.Handlers) target.Dispatch
)

// CheckAstral is a default implementation of Astral function. Returns error if astral is not started.
func CheckAstral(_ context.Context) error { return apphost.Check() }

func NewFeat(
	astral Astral,
	port target.Port,
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
		astral: astral,
		port:   port,
		tray:   tray,
		serve:  service(handlers),
	}
}

// Run portal service including astral daemon if not started. Optionally displays an indicator in OS tray.
func (f Feat) Run(
	ctx context.Context,
	tray bool,
) (err error) {
	if err = f.astral(ctx); err != nil {
		return
	}

	log := plog.Get(ctx).Type(f).Set(&ctx)
	request := rpc.NewRequest(id.Anyone, f.port.String())
	if err = rpc.Command(request, "ping"); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
		return
	}
	err = nil
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		if err := f.serve(ctx, f.port.String()); err != nil {
			log.Printf("serve exit: %v\n", err)
		} else {
			log.Println("serve exit")
		}
	}()
	if tray && f.tray != nil {
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
