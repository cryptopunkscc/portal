package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	astral "github.com/cryptopunkscc/portal/runtime/apphost"
	api "github.com/cryptopunkscc/portal/runtime/portal"
	"github.com/cryptopunkscc/portal/runtime/rpc"
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
func CheckAstral(_ context.Context) error { return astral.Check() }

type Deps interface {
	Port() target.Port
	Open() target.Request
	Astral() Astral
	Handlers() Handlers
	Observe() Observe
	Shutdown() context.CancelFunc
	Client() apphost.Client
}

func Feat(d Deps) target.Request {
	astral := d.Astral()
	client := d.Client()
	port := d.Port()
	handlers := d.Handlers()
	maps.Copy(handlers, Handlers{
		"ping":    func() {},
		"open":    d.Open(),
		"observe": d.Observe(),
		"close":   d.Shutdown(),
	})
	return func(ctx context.Context, src string, args ...string) (err error) {
		if err = astral(ctx); err != nil {
			return plog.Err(err)
		}
		if err = check(port); err != nil {
			return plog.Err(err)
		}
		if err = serve(ctx, client, port, handlers); err != nil {
			return plog.Err(err)
		}
		return
	}
}

func check(port target.Port) (err error) {
	if err = api.Client(port.String()).Ping(); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
	}
	return nil
}

func serve(
	ctx context.Context,
	client apphost.Client,
	port target.Port,
	handlers Handlers,
) (err error) {
	log := plog.Get(ctx)
	log.Printf("serve start at port:%s", port)
	defer log.Printf("serve exit:%s", port)
	if err = rpc.NewApp(port.String()).
		Client(client).
		Routes("*").
		RouteMap(rpc.Handlers(handlers)).
		Run(ctx); err != nil {
		log.Printf("serve error: %v\n", err)
		return plog.Err(err)
	}
	return
}
