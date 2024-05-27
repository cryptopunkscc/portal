package builder

import (
	"context"
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/runner/spawn"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"log"
	"reflect"
	"sync"
)

type Scope[T target.Portal] struct {
	Prefix      []string
	Port        string
	WaitGroup   *sync.WaitGroup
	TargetCache *target.Cache[T]
	RpcHandlers rpc.Handlers

	WrapApi      func(target.Api) target.Api
	NewTargetRun func(target.NewApi) target.Run[T]
	NewTray      func(target.Dispatch) target.Tray
	NewServe     func(rpc.Handlers) target.Dispatch

	TargetFinder target.Finder[T]
	ExecTarget   target.Run[T]

	AppsPath        target.Path
	DispatchTarget  target.Dispatch
	DispatchService target.Dispatch

	FeatObserve   func(ctx context.Context, conn rpc.Conn) (err error)
	FeatInstall   func(src string) error
	FeatUninstall func(src string) error

	// auto dependencies
	TargetFind   target.Find[T]
	FeatDispatch target.Dispatch
	FeatOpen     target.Dispatch
	FeatServe    func(context.Context, bool) error
	FeatList     func() []target.App
}

func (s *Scope[T]) GetWait() *sync.WaitGroup            { return assert(s.WaitGroup) }
func (s *Scope[T]) GetExecTarget() target.Run[T]        { return assert(s.ExecTarget) }
func (s *Scope[T]) GetTargetFinder() target.Finder[T]   { return assert(s.TargetFinder) }
func (s *Scope[T]) GetTargetCache() *target.Cache[T]    { return assert(s.TargetCache) }
func (s *Scope[T]) GetDispatchTarget() target.Dispatch  { return assert(s.DispatchTarget) }
func (s *Scope[T]) GetDispatchService() target.Dispatch { return assert(s.DispatchService) }

func (s *Scope[T]) GetTargetFind() target.Find[T] {
	if s.TargetFind == nil {
		launcherSvelteFs := embedApps.LauncherSvelteFS
		resolveEmbed := portal.NewResolver[target.App](
			apps.Resolve[target.App](),
			source.FromFS(launcherSvelteFs),
		)
		findPath := target.Mapper[string, string](
			resolveEmbed.Path,
			assert(s.AppsPath),
		)
		s.TargetFind = s.GetTargetFinder().Cached(s.GetTargetCache())(findPath, launcherSvelteFs)
	}
	return s.TargetFind
}

func (s *Scope[T]) GetServeFeature() func(context.Context, bool) error {
	runSpawn := spawn.NewRunner(s.GetWait(), s.GetTargetFind(), s.GetExecTarget()).Run
	runTray := target.Tray(nil)
	if s.NewTray != nil {
		runTray = s.NewTray(runSpawn)
	}
	if s.RpcHandlers == nil {
		s.RpcHandlers = rpc.Handlers{}
	}
	return serve.NewFeat(
		assert(s.Port),
		assert(s.NewServe),
		assert(s.RpcHandlers),
		assert(runSpawn),
		assert(s.FeatObserve),
		runTray,
	)
}

func (s *Scope[T]) GetOpenFeature() target.Dispatch {
	if s.FeatOpen == nil {
		newApphost := apphost.NewFactory(s.GetDispatchTarget(), s.Prefix...)
		newApi := target.ApiFactory(assert(s.WrapApi),
			newApphost.NewAdapter,
			newApphost.WithTimeout,
		)
		s.FeatOpen = open.NewFeat[T](s.GetTargetFind(), s.NewTargetRun(newApi))
	}
	return s.FeatOpen
}

func (s *Scope[T]) GetDispatchFeature() target.Dispatch {
	if s.FeatDispatch == nil {
		s.FeatDispatch = dispatch.NewFeat(s.GetDispatchTarget(), s.GetDispatchService())
	}
	return s.FeatDispatch
}

func assert[T any](arg T) T {
	check(arg)
	return arg
}

func check(arg any) {
	if arg == nil || arg == "" || reflect.ValueOf(arg).IsZero() {
		log.Panicf("nil dependency: %T", arg)
	}
}
