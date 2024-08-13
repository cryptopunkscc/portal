package main

import (
	"context"
	manifest "github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/clir"
	"github.com/cryptopunkscc/portal/factory/srv"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	signal "github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runtime/msg"
	. "github.com/cryptopunkscc/portal/target"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	mod := deps[Portal_]{}
	mod.Deps = &mod
	mod.CancelFunc = cancel
	go signal.OnShutdown(cancel)
	println("...")
	plog.ErrorStackTrace = true
	log := plog.New().D().Scope("dev").Set(&ctx)
	log.Println("starting portal development", os.Args)
	defer log.Println("closing portal development")
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx, manifest.NameDev, manifest.DescriptionDev, version.Run)
	cli.Dev(start.Feat(&mod))
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
	mod.WaitGroup().Wait()
}

type deps[T Portal_] struct{ srv.Module[T] }

func (d *deps[T]) Executable() string   { return "portal-dev" }
func (d *deps[T]) Request() Request     { return query.Request.Run }
func (d *deps[T]) Serve() Request       { return serve.Feat(d).Start }
func (d *deps[T]) Astral() serve.Astral { return serve.CheckAstral }
func (d *deps[T]) Handlers() serve.Handlers {
	return serve.Handlers{
		PortMsg.Name: msg.NewBroadcast(PortMsg, d.Processes()).BroadcastMsg,
	}
}
func (d *deps[T]) Resolve() Resolve[T] { return sources.Resolver[T]() }
func (d *deps[T]) Run() Run[T] {
	return multi.Runner[T](
		app.Run(exec.Portal[PortalJs]("portal-dev-goja", "o").Run),
		app.Run(exec.Portal[PortalHtml]("portal-dev-wails", "o").Run),
		app.Run(exec.Portal[ProjectGo]("portal-dev-go", "o").Run),
		app.Run(exec.Portal[AppExec]("portal-dev-exec", "o").Run),
	)
}
func (d *deps[T]) Priority() Priority {
	return Priority{
		Match[Project_],
		Match[Dist_],
		Match[Bundle_],
	}
}
