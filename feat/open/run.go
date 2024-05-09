package open

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"log"
	"time"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error) {
	log.Println("open.Run", src, attach)
	// dispatch execution to service
	if !attach {
		return Serve(ctx, bindings, src)
	}

	// resolve apps from given source
	apps, err := portal.ResolveApps(src)
	if len(apps) == 0 {
		return errors.Join(fmt.Errorf("no apps found in %s", src), err)
	}

	// execute multiple targets as separate processes
	if len(apps) > 1 {
		return portal.Spawn(ctx, apps)
	}

	// execute single target in current process
	for _, app := range apps {
		_ = portal.Attach(ctx, bindings, app)
	}

	return
}

func Serve(
	ctx context.Context,
	bindings runtime.New,
	src string,
) (err error) {
	// dispatch query to service
	if err = portal.SrvOpenCtx(ctx, src); err == nil {
		return
	}

	// wait up to 2 seconds for service start and execute given source as child process
	go func() {
		if err = portal.Await(2 * time.Second); err != nil {
			log.Println("serve await timout:", err)
		}
		go func() {
			<-ctx.Done()
			log.Println("serve ctx done")
		}()
		if err = portal.SrvOpenCtx(ctx, src); err != nil {
			log.Printf("serve portal.open %s: %v", src, err)
		}
	}()

	// start portal service
	if err = serve.Run(ctx, bindings, true); err == nil {
		err = fmt.Errorf("cannot serve portal: %w", err)
	}
	return
}
