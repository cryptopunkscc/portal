package serve

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func Runner(d Deps) target.Run[string] {
	startAstral := d.Astral()
	handler := Handler(d)
	return func(ctx context.Context, _ string, _ ...string) (err error) {
		if err = startAstral(ctx); err != nil {
			return plog.Err(err)
		}
		if err = checkPortald(); err != nil {
			return plog.Err(err)
		}
		if err = serve(ctx, handler); err != nil {
			return plog.Err(err)
		}
		return
	}
}

type Deps interface {
	Service
	Astral() Astral
}

// Astral starts daemon if not already running.
type Astral func(ctx context.Context) (err error)

func checkPortald() (err error) {
	if err = portal.DefaultClient.Ping(); err == nil {
		err = fmt.Errorf("port already registered or astral is not running: %v", err)
	}
	return nil
}

func serve(
	ctx context.Context,
	handler cmd.Handler,
) (err error) {
	log := plog.Get(ctx).Scope("serve")
	log.Println("serve start")
	defer log.Printf("serve exit")

	router := apphost.Default().Router(handler)
	router.Logger = log
	err = router.Run(ctx)

	if err != nil {
		log.Printf("serve error: %v\n", err)
		return plog.Err(err)
	}
	return
}
