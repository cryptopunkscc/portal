package main

import (
	"context"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/di/srv"
	"github.com/cryptopunkscc/portal/dispatch/query"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	signal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runtime/msg"
	. "github.com/cryptopunkscc/portal/target"
	"os"
)

func main() {
	mod := Module[Portal_]{}
	mod.Deps = &mod
	ctx, cancel := context.WithCancel(context.Background())
	mod.CancelFunc = cancel
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(mod.Start())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type Module[T Portal_] struct{ srv.Module[T] }

func (d *Module[T]) Executable() string        { return "portal-dev" }
func (d *Module[T]) Start() Dispatch           { return start.Inject(d).Run }
func (d *Module[T]) JoinTarget() Dispatch      { return query.NewOpen().Run }
func (d *Module[T]) DispatchService() Dispatch { return serve.Dispatch(d) }
func (d *Module[T]) Astral() serve.Astral      { return serve.CheckAstral }
func (d *Module[T]) Handlers() serve.Handlers {
	return serve.Handlers{
		PortMsg.Name: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg,
	}
}
func (d *Module[T]) Resolve() Resolve[T] { return sources.Resolver[T]() }
func (d *Module[T]) Run() Run[T] {
	return multi.NewRunner[T](
		app.Run(exec.NewPortal[PortalJs]("portal-dev-goja", "o").Run),
		app.Run(exec.NewPortal[PortalHtml]("portal-dev-wails", "o").Run),
		app.Run(exec.NewPortal[ProjectGo]("portal-dev-go", "o").Run),
		app.Run(exec.NewPortal[AppExec]("portal-dev-exec", "o").Run),
	).Run
}
func (d *Module[T]) Priority() Priority {
	return Priority{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
