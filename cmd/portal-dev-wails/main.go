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
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/wails"
	"github.com/cryptopunkscc/portal/runner/wails_dev"
	"github.com/cryptopunkscc/portal/runner/wails_dist"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
	"github.com/cryptopunkscc/portal/target/find"
	"github.com/cryptopunkscc/portal/target/portals"
)

func main() {
	deps := NewDeps()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-wails").Set(&ctx)
	go sig.OnShutdown(cancel)
	portalPort.InitPrefix("dev")
	cli := clir.NewCli(ctx,
		"Portal-dev-wails",
		"Portal html development driven by wails.",
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
	cache *Cache[PortalHtml]
}
type Adapter struct{ Api }

func NewDeps() *Deps                             { return &Deps{cache: NewCache[PortalHtml]()} }
func NewAdapter(api Api) Api                     { return &Adapter{Api: api} }
func (d *Deps) Path() Path                       { return featApps.Path }
func (d *Deps) NewApi() NewApi                   { return di.Single(api.New, d) }
func (d *Deps) WrapApi(api Api) Api              { return NewAdapter(api) }
func (d *Deps) TargetRun() Run[PortalHtml]       { return d.NewRunTarget(d.NewApi()) }
func (d *Deps) TargetFind() Find[PortalHtml]     { return di.Single(find.New[PortalHtml], d) }
func (d *Deps) TargetFinder() Finder[PortalHtml] { return portals.NewFind[PortalHtml] }
func (d *Deps) TargetCache() *Cache[PortalHtml]  { return d.cache }
func (d *Deps) TargetDispatch() Dispatch         { return query.NewRunner[PortalHtml](PortOpen).Start }
func (d *Deps) FeatOpen() Dispatch               { return di.Single(open.Inject[PortalHtml], d) }
func (d *Deps) NewRunTarget(newApi NewApi) Run[PortalHtml] {
	return multi.NewRunner[PortalHtml](
		reload.Immutable(newApi, PortMsg, wails_dev.NewRunner), // FIXME propagate sendMsg
		reload.Mutable(newApi, PortMsg, wails_dist.NewRunner),
		reload.Immutable(newApi, PortMsg, wails.NewRunner),
	).Run
}
