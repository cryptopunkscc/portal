package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	api "github.com/cryptopunkscc/portal/runtime/portal"
	apphost2 "github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func Runner(d Deps) target.Run[string] {
	startAstral := d.Astral()
	port := d.Port()
	handler := Handler(d)
	return func(ctx context.Context, _ string, _ ...string) (err error) {
		if err = startAstral(ctx); err != nil {
			return plog.Err(err)
		}
		if err = checkPortald(port); err != nil {
			return plog.Err(err)
		}
		if err = serve(ctx, port, handler); err != nil {
			return plog.Err(err)
		}
		return
	}
}

type Deps interface {
	Service
	Astral() Astral
	Port() apphost.Port
}

// Astral starts daemon if not already running.
type Astral func(ctx context.Context) (err error)

func checkPortald(port apphost.Port) (err error) {
	if err = api.Client(port.String()).Ping(); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
	}
	return nil
}

func serve(
	ctx context.Context,
	port apphost.Port,
	handler cmd.Handler,
) (err error) {
	log := plog.Get(ctx).Scope("serve")
	log.Printf("serve start at port:%s", port)
	defer log.Printf("serve exit:%s", port)

	router := apphost2.NewRouter(handler, port)
	router.Logger = log
	err = router.Run(ctx)

	if err != nil {
		log.Printf("serve error: %v\n", err)
		return plog.Err(err)
	}
	return
}
