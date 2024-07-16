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
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/query"
	"github.com/cryptopunkscc/portal/runner/reload"
	. "github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/api"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/cache"
	"github.com/cryptopunkscc/portal/target/find"
)

func main() {
	deps := NewDeps[AppExec]()
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("dev-exec").Set(&ctx)
	go sig.OnShutdown(cancel)
	cli := clir.NewCli(ctx,
		"Portal-dev-exec",
		"Portal js development runner for executables.",
		version.Run,
	)
	cli.Open(deps.FeatOpen())
	if err := cli.Run(); err != nil {
		log.Println(err)
	}
	cancel()
}

type Deps[T AppExec] struct {
	di.Cache
	cache *Cache[T]
}
type Adapter struct{ Api }

func NewDeps[T AppExec]() *Deps[T]          { return &Deps[T]{} }
func NewAdapter(api Api) Api                { return &Adapter{Api: api} }
func (d *Deps[T]) Path() Path               { return featApps.Path }
func (d *Deps[T]) Executable() string       { return "portal-dev" }
func (d *Deps[T]) GetCacheDir() string      { return di.Single(cache.Dir, d) }
func (d *Deps[T]) NewApi() NewApi           { return di.Single(api.New, d) }
func (d *Deps[T]) WrapApi(api Api) Api      { return NewAdapter(api) }
func (d *Deps[T]) TargetRun() Run[T]        { return d.NewRunTarget(d.NewApi()) }
func (d *Deps[T]) TargetCache() *Cache[T]   { return d.cache }
func (d *Deps[T]) TargetFind() Find[T]      { return di.Single(find.New[T], d) }
func (d *Deps[T]) TargetFinder() Finder[T]  { return apps.NewFind[T] }
func (d *Deps[T]) TargetDispatch() Dispatch { return query.NewRunner[T](PortOpen).Start }
func (d *Deps[T]) FeatOpen() Dispatch       { return di.Single(open.Inject[T], d) }
func (d *Deps[T]) NewRunTarget(newApi NewApi) Run[T] {
	return multi.NewRunner[T](
		reload.Immutable(newApi, PortMsg, reload.Adapter(exec.NewBundleRunner(d.GetCacheDir()))),
		reload.Immutable(newApi, PortMsg, reload.Adapter(exec.NewDistRunner())),
	).Run
}
