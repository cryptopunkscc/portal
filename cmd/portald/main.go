package main

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/api/portal"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/request"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	singal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	runtime "github.com/cryptopunkscc/portal/runtime/portal"
	apphost2 "github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
	"sync"
)

func main() {
	d := &deps[Portal_]{}
	d.check()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("portald").Set(&ctx)
	d.CancelFunc = cancel
	go singal.OnShutdown(cancel)
	if err := d.Router().Run(ctx); err != nil {
		log.Println(err)
	}
	cancel()
	d.WaitGroup().Wait()
}

type deps[T Portal_] struct {
	CancelFunc context.CancelFunc
	wg         sync.WaitGroup
	processes  sig.Map[string, T]
}

func (d *deps[T]) check() {
	if err := apphost.Check(); err != nil {
		panic(err)
	}
	if err := runtime.Client(d.Port().String()).Ping(); err == nil {
		err = fmt.Errorf("port already registered or astral not running: %v", err)
		panic(err)
	}
}
func (d *deps[T]) Executable() string             { return "portal" }
func (d *deps[T]) WaitGroup() *sync.WaitGroup     { return &d.wg }
func (d *deps[T]) Processes() *sig.Map[string, T] { return &d.processes }
func (d *deps[T]) Shutdown() context.CancelFunc   { return d.CancelFunc }
func (d *deps[T]) Resolve() Resolve[T]            { return sources.Resolver[T]() }
func (d *deps[T]) Open() Request                  { return request.Create[T](d) }
func (d *deps[T]) Port() port.Port                { return port.New(d.Executable()) }
func (d *deps[T]) Handler() cmd.Handler           { return portal.Handler(d) }
func (d *deps[T]) Router() *apphost2.Router       { return apphost2.NewRouter(d.Handler(), d.Port()) }
func (d *deps[T]) Run() Run[T] {
	return multi.Runner[T](
		app.Run(exec.Portal[AppJs]("portal-app-goja", "o").Run),
		app.Run(exec.Portal[AppHtml]("portal-app-wails", "o").Run),
		app.Run(exec.Bundle(CacheDir(d.Executable())).Run),
	)
}
func (d *deps[T]) Priority() Priority {
	return []Matcher{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
