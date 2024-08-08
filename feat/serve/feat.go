package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
	"maps"
)

type (
	// Astral starts daemon if not already running.
	Astral func(ctx context.Context) (err error)

	// Observe on installed applications.
	Observe func(ctx context.Context, conn rpc.Conn) (err error)

	Handlers map[string]any
)

// CheckAstral is a default implementation of Astral function. Returns error if astral is not started.
func CheckAstral(_ context.Context) error { return apphost.Check() }

type Deps interface {
	Port() target.Port
	Open() target.Dispatch
	Astral() Astral
	Handlers() Handlers
	Observe() Observe
	Shutdown() context.CancelFunc
}

func Dispatch(d Deps) target.Dispatch {
	run := Run(d)
	return func(ctx context.Context, _ string, _ ...string) (err error) {
		go func() {
			if err = run(ctx); err != nil {
				plog.Get(ctx).Type(d).Println("dispatch:", err)
			}
		}()
		return
	}
}

func Run(d Deps) func(ctx context.Context) error {
	astral := d.Astral()
	port := d.Port()
	handlers := d.Handlers()
	maps.Copy(handlers, Handlers{
		"ping":    func() {},
		"open":    d.Open(),
		"observe": d.Observe(),
		"close":   d.Shutdown(),
	})
	return func(ctx context.Context) (err error) {
		if err = astral(ctx); err != nil {
			return plog.Err(err)
		}
		if err = check(port); err != nil {
			return plog.Err(err)
		}
		if err = serve(ctx, port, handlers); err != nil {
			return plog.Err(err)
		}
		return
	}
}

func check(port target.Port) (err error) {
	request := rpc.NewRequest(id.Anyone, port.String())
	if err = rpc.Command(request, "ping"); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
	}
	return nil
}

func serve(
	ctx context.Context,
	port target.Port,
	handlers Handlers,
) (err error) {
	log := plog.Get(ctx).Set(&ctx)
	log.Printf("serve start at port:%s", port)
	defer log.Printf("serve exit:%s", port)
	if err = rpc.NewApp(port.String()).
		Routes("*").
		RouteMap(rpc.Handlers(handlers)).
		Run(ctx); err != nil {
		log.Printf("serve error: %v\n", err)
		return plog.Err(err)
	}
	return
}
