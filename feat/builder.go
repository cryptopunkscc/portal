package feat

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	embedApps "github.com/cryptopunkscc/go-astral-js/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/dispatch"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
	"github.com/cryptopunkscc/go-astral-js/target/apphost"
	"github.com/cryptopunkscc/go-astral-js/target/apps"
	"github.com/cryptopunkscc/go-astral-js/target/portal"
	"github.com/cryptopunkscc/go-astral-js/target/source"
	"github.com/cryptopunkscc/go-astral-js/target/spawn"
	"log"
	"reflect"
	"sync"
)

type Scope[T target.Portal] struct {
	Prefix      []string
	Port        target.Port
	WaitGroup   *sync.WaitGroup
	TargetCache *target.Cache[T]
	RpcHandlers rpc.Handlers
	Processes   *sig.Map[string, T]

	WrapApi       func(target.Api) target.Api
	NewRunTarget  func(target.NewApi) target.Run[T]
	NewRunTray    func(target.Dispatch) target.Tray
	NewRunService func(rpc.Handlers) target.Dispatch

	ExecTarget   target.Run[T]
	TargetFinder target.Finder[T]

	GetPath         target.Path
	DispatchTarget  target.Dispatch
	DispatchService target.Dispatch
	JoinTarget      target.Dispatch

	FeatObserve func(ctx context.Context, conn rpc.Conn) (err error)

	// auto dependencies
	TargetFind   target.Find[T]
	FeatDispatch target.Dispatch
	FeatOpen     target.Dispatch
	FeatServe    *serve.Feat
	FeatList     func() []target.App
}

func (s *Scope[T]) GetWaitGroup() *sync.WaitGroup       { return assert(s.WaitGroup) }
func (s *Scope[T]) GetProcesses() *sig.Map[string, T]   { return assert(s.Processes) }
func (s *Scope[T]) GetExecTarget() target.Run[T]        { return assert(s.ExecTarget) }
func (s *Scope[T]) GetTargetFinder() target.Finder[T]   { return assert(s.TargetFinder) }
func (s *Scope[T]) GetTargetCache() *target.Cache[T]    { return assert(s.TargetCache) }
func (s *Scope[T]) GetJoinTarget() target.Dispatch      { return assert(s.JoinTarget) }
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
			assert(s.GetPath),
		)
		s.TargetFind = s.GetTargetFinder().Cached(s.GetTargetCache())(findPath, launcherSvelteFs)
	}
	return s.TargetFind
}

func (s *Scope[T]) GetServeFeature() *serve.Feat {
	if s.FeatServe == nil {
		runSpawn := spawn.NewRunner(
			s.GetWaitGroup(),
			s.GetTargetFind(),
			s.GetExecTarget(),
			s.GetProcesses(),
		).Run
		runTray := target.Tray(nil)
		if s.NewRunTray != nil {
			runTray = s.NewRunTray(runSpawn)
		}

		s.FeatServe = serve.NewFeat(
			assert(s.Port),
			assert(s.NewRunService),
			s.RpcHandlers,
			assert(runSpawn),
			assert(s.FeatObserve),
			runTray,
		)
	}
	return s.FeatServe
}

func (s *Scope[T]) GetOpenFeature() target.Dispatch {
	if s.FeatOpen == nil {
		newApphost := apphost.NewFactory(s.GetDispatchTarget(), s.Port.Prefix...)
		newApi := target.ApiFactory(assert(s.WrapApi),
			newApphost.NewAdapter,
			newApphost.WithTimeout,
		)
		s.FeatOpen = open.NewFeat[T](s.GetTargetFind(), s.NewRunTarget(newApi))
	}
	return s.FeatOpen
}

func (s *Scope[T]) GetDispatchFeature() target.Dispatch {
	if s.FeatDispatch == nil {
		s.FeatDispatch = dispatch.NewFeat(s.GetJoinTarget(), s.GetDispatchService())
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
