package portal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"time"
)

type Runner[T target.Portal] struct {
	runtime.New
	runtime.Tray
	runtime.Serve
	runtime.Attach[T]
	runtime.Resolve[T]
	rpc.Handlers
	Cmd string
}

func (r Runner[T]) Run(
	ctx context.Context,
	src string,
	attach bool,
) (err error) {

	if attach {
		return r.attach(ctx, src)
	} else {
		return r.dispatch(ctx, src)
	}
}

func (r Runner[T]) attach(ctx context.Context, src string) (err error) {

	// resolve apps from given source
	apps, err := r.Resolve(src)
	if len(apps) == 0 {
		return errors.Join(fmt.Errorf("no apps found in %s", src), err)
	}

	// execute multiple targets as separate processes
	if len(apps) > 1 {
		return Spawn(nil, ctx, apps, r.Cmd)
	}

	// execute single target in current process
	for _, app := range apps {
		m := app.Manifest()
		log.Printf("running %s %s %s %s", m.Name, m.Version, m.Package, app.Path())
		_ = r.Attach(ctx, r.New, app)
	}
	return
}

func (r Runner[T]) dispatch(
	ctx context.Context,
	src string,
) (err error) {
	// dispatch query to service
	if err = SrvOpenCtx(ctx, src); err == nil {
		return
	}

	// wait up to 2 seconds for service start and execute given source as child process
	go func() {
		if err = Await(2 * time.Second); err != nil {
			log.Println("serve await timout:", err)
		}
		go func() {
			<-ctx.Done()
			log.Println("serve ctx done")
		}()
		if err = SrvOpenCtx(ctx, src); err != nil {
			log.Printf("serve portal.open %s: %v", src, err)
		}
	}()

	// start portal service
	if err = r.Serve(ctx, r.New, r.Handlers, r.Tray); err == nil {
		err = fmt.Errorf("portal serve exit: %v", err)
	}
	return
}
