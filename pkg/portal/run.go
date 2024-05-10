package portal

import (
	"context"
	"errors"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"strings"
	"time"
)

type Runner[T target.Portal] struct {
	runtime.New
	runtime.Tray
	runtime.Serve
	runtime.Attach[T]
	runtime.Resolve[T]
	rpc.Handlers
	Action  string
	Port    string
	Request rpc.Conn
}

func (r Runner[T]) Run(
	ctx context.Context,
	src string,
	attach bool,
) (err error) {
	r.Request = rpc.NewRequest(id.Anyone, r.Port)
	r.Request.Logger(log.New(log.Writer(), r.Port+" ", 0))
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
		return Spawn(nil, ctx, apps, r.Action)
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

	chunks := []string{src}
	if strings.HasPrefix(r.Port, "dev") {
		chunks = []string{"dev", src}
	}

	// dispatch query to service
	if err = SrvOpenCtx(ctx, chunks...); err == nil {
		return
	}

	// wait up to 2 seconds for service start and execute given source as child process
	go func(chunks []string) {

		if err = r.Await(2 * time.Second); err != nil {
			log.Println("serve await timout:", err)
		}
		go func() {
			<-ctx.Done()
		}()
		if err = SrvOpenCtx(ctx, chunks...); err != nil {
			log.Printf("serve portal.open %s: %v", src, err)
		}
	}(chunks)

	// start portal service
	if err = r.Serve(ctx, r.Port, r.Handlers, r.Tray); err == nil {
		err = fmt.Errorf("portal serve exit: %v", err)
	}
	return
}
