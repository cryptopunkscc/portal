package main

import (
	"context"
	"github.com/cryptopunkscc/portal/clir"
	featApps "github.com/cryptopunkscc/portal/feat/apps"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/di"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/wails"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/find"
)

func main() {
	deps := NewDeps()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app-wails").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-wails",
		"Portal html runner driven by wails.",
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
	cache *Cache[AppHtml]
}

type Adapter struct{ Api }

func NewDeps() *Deps                                    { return &Deps{cache: NewCache[AppHtml]()} }
func NewAdapter(api Api) Api                            { return &Adapter{Api: api} }
func (d *Deps) Path() Path                              { return featApps.Path }
func (d *Deps) NewApi() NewApi                          { return di.Single(api.New, d) }
func (d *Deps) WrapApi(api Api) Api                     { return NewAdapter(api) }
func (d *Deps) TargetRun() Run[AppHtml]                 { return d.NewRunTarget(d.NewApi()) }
func (d *Deps) TargetFind() Find[AppHtml]               { return di.Single(find.New[AppHtml], d) }
func (d *Deps) TargetFinder() Finder[AppHtml]           { return apps.NewFind[AppHtml] }
func (d *Deps) TargetCache() *Cache[AppHtml]            { return d.cache }
func (d *Deps) TargetDispatch() Dispatch                { return query.NewRunner[AppHtml](PortOpen).Start }
func (d *Deps) NewRunTarget(newApi NewApi) Run[AppHtml] { return wails.NewRun(newApi) }
func (d *Deps) FeatOpen() Dispatch                      { return open.Inject[AppHtml](d) }
