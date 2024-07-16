package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	portalPort "github.com/cryptopunkscc/portal/pkg/port"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/runner/goja_dev"
	"github.com/cryptopunkscc/portal/runner/goja_dist"
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
	log := plog.New().D().Scope("dev-goja").Set(&ctx)
	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")
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
	cache *Cache[PortalJs]
}
type Adapter struct{ Api }

func NewDeps() *Deps                           { return &Deps{cache: NewCache[PortalJs]()} }
func NewAdapter(api Api) Api                   { return &Adapter{Api: api} }
func (d *Deps) Path() Path                     { return featApps.Path }
func (d *Deps) NewApi() NewApi                 { return di.Single(api.New, d) }
func (d *Deps) WrapApi(api Api) Api            { return NewAdapter(api) }
func (d *Deps) TargetRun() Run[PortalJs]       { return d.NewRunTarget(d.NewApi()) }
func (d *Deps) TargetFind() Find[PortalJs]     { return di.Single(find.New[PortalJs], d) }
func (d *Deps) TargetFinder() Finder[PortalJs] { return portals.NewFind[PortalJs] }
func (d *Deps) TargetCache() *Cache[PortalJs]  { return d.cache }
func (d *Deps) TargetDispatch() Dispatch       { return query.NewRunner[PortalJs](PortOpen).Start }
func (d *Deps) FeatOpen() Dispatch             { return di.Single(open.Inject[PortalJs], d) }
func (d *Deps) NewRunTarget(newApi NewApi) Run[PortalJs] {
	return multi.NewRunner[PortalJs](
		reload.Mutable(newApi, PortMsg, goja_dev.NewRunner),
		reload.Mutable(newApi, PortMsg, goja_dist.NewRunner),
		reload.Immutable(newApi, PortMsg, goja.NewRunner),
	).Run
}
