package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/cryptopunkscc/portal/target/portals"
)

func main() {
	deps := NewDeps()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-go").Set(&ctx)
	go sig.OnShutdown(cancel)
	port.InitPrefix("dev")
	cli := clir.NewCli(ctx,
		"Portal-dev-goja",
		"Portal js development driven by goja.",
		version.Run,
	)
	cli.Open(deps.FeatOpen())

	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Deps struct {
	di.Cache
	cache *Cache[ProjectGo]
}
type Adapter struct{ Api }

func NewDeps() *Deps                            { return &Deps{cache: NewCache[ProjectGo]()} }
func NewAdapter(api Api) Api                    { return &Adapter{Api: api} }
func (d *Deps) Path() Path                      { return featApps.Path }
func (d *Deps) NewApi() NewApi                  { return di.Single(api.New, d) }
func (d *Deps) WrapApi(api Api) Api             { return NewAdapter(api) }
func (d *Deps) TargetRun() Run[ProjectGo]       { return d.NewRunTarget(d.NewApi()) }
func (d *Deps) TargetFind() Find[ProjectGo]     { return di.Single(find.New[ProjectGo], d) }
func (d *Deps) TargetFinder() Finder[ProjectGo] { return portals.NewFind[ProjectGo] }
func (d *Deps) TargetCache() *Cache[ProjectGo]  { return d.cache }
func (d *Deps) TargetDispatch() Dispatch        { return query.NewRunner[ProjectGo](PortOpen).Start }
func (d *Deps) FeatOpen() Dispatch              { return di.Single(open.Inject[ProjectGo], d) }
func (d *Deps) NewRunTarget(newApi NewApi) Run[ProjectGo] {
	var runner = go_dev.NewAdapter(app.Run(exec.NewDistRunner().Run))
	return multi.NewRunner[ProjectGo](reload.Mutable(newApi, PortMsg, runner)).Run
}
