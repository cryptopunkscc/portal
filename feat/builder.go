package feat

import (
	"context"
	"github.com/cryptopunkscc/astrald/sig"
	embedApps "github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/feat/dispatch"
	"github.com/cryptopunkscc/portal/feat/open"
	"github.com/cryptopunkscc/portal/feat/serve"
	"github.com/cryptopunkscc/portal/pkg/assets"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/apphost"
	"github.com/cryptopunkscc/portal/target/apps"
	"github.com/cryptopunkscc/portal/target/portal"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/cryptopunkscc/portal/target/spawn"
	"log"
	"os"
	"path"
	"reflect"
	"sync"
)

type Scope[T target.Portal] struct {
	CacheDir    string
	Executable  string
	Port        target.Port
	NewApi      target.NewApi
	TargetCache *target.Cache[T]
	RpcHandlers rpc.Handlers
	Processes   *sig.Map[string, T]

	Astral        serve.Astral
	WrapApi       func(target.Api) target.Api
	NewRunTarget  func(target.NewApi) target.Run[T]
	NewRunTray    func(target.Dispatch) target.Tray
	NewRunService func(rpc.Handlers) target.Dispatch
	NewExecTarget func(string, string) target.Run[T] // TODO replace factory with direct reference

	TargetFinder    target.Finder[T]
	GetPath         target.Path
	DispatchTarget  target.Dispatch
	DispatchService target.Dispatch
	JoinTarget      target.Dispatch

	FeatObserve func(ctx context.Context, conn rpc.Conn) (err error)

	// auto dependencies
	WaitGroup    sync.WaitGroup
	TargetFind   target.Find[T]
	FeatDispatch target.Dispatch
	FeatOpen     target.Dispatch
	FeatServe    *serve.Feat
	FeatList     func() []target.App
}

func (s *Scope[T]) GetPort() target.Port                { return require(s.Port) }
func (s *Scope[T]) GetExecutable() string               { return require(s.Executable) }
func (s *Scope[T]) GetWaitGroup() *sync.WaitGroup       { return &s.WaitGroup }
func (s *Scope[T]) GetProcesses() *sig.Map[string, T]   { return require(s.Processes) }
func (s *Scope[T]) GetTargetFinder() target.Finder[T]   { return require(s.TargetFinder) }
func (s *Scope[T]) GetTargetCache() *target.Cache[T]    { return require(s.TargetCache) }
func (s *Scope[T]) GetJoinTarget() target.Dispatch      { return require(s.JoinTarget) }
func (s *Scope[T]) GetDispatchTarget() target.Dispatch  { return require(s.DispatchTarget) }
func (s *Scope[T]) GetDispatchService() target.Dispatch { return require(s.DispatchService) }

func (s *Scope[T]) GetExecTarget() target.Run[T] {
	return require(s.NewExecTarget)(
		s.GetCacheDir(),
		s.GetExecutable(),
	)
}

func (s *Scope[T]) GetCacheDir() string {
	if s.CacheDir == "" {
		var err error
		if s.CacheDir, err = os.UserCacheDir(); err != nil {
			panic(err)
		}
		s.CacheDir = path.Join(s.CacheDir, s.Executable)
	}
	return s.CacheDir
}

func (s *Scope[T]) GetNewApi() target.NewApi {
	if s.NewApi == nil {
		apphost.ConnectionsThreshold = 0
		newApphost := apphost.NewFactory(s.GetDispatchTarget())
		s.NewApi = target.ApiFactory(require(s.WrapApi),
			newApphost.NewAdapter,
			newApphost.WithTimeout,
		)
	}
	return s.NewApi
}

func (s *Scope[T]) GetTargetFind() target.Find[T] {
	if s.TargetFind == nil {
		launcherSvelteFs := embedApps.LauncherSvelteFS
		resolveEmbedApp := portal.NewResolver[target.App](
			apps.Resolve[target.App](),
			source.FromFS(launcherSvelteFs),
		)
		findPath := target.Mapper[string, string](
			resolveEmbedApp.Path,
			require(s.GetPath),
		)
		s.TargetFind = s.GetTargetFinder().Cached(s.GetTargetCache())(
			findPath,
			launcherSvelteFs,
			assets.OsFS{},
		)
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

		var runTray target.Tray
		if s.NewRunTray != nil {
			runTray = s.NewRunTray(runSpawn)
		}

		s.FeatServe = serve.NewFeat(
			require(s.Astral),
			s.GetPort(),
			require(s.NewRunService),
			s.RpcHandlers,
			require(runSpawn),
			require(s.FeatObserve),
			runTray,
		)
	}
	return s.FeatServe
}

func (s *Scope[T]) GetOpenFeature() target.Dispatch {
	if s.FeatOpen == nil {
		s.FeatOpen = open.NewFeat[T](
			s.GetTargetFind(),
			s.NewRunTarget(s.GetNewApi()),
		)
	}
	return s.FeatOpen
}

func (s *Scope[T]) GetDispatchFeature() target.Dispatch {
	if s.FeatDispatch == nil {
		s.FeatDispatch = dispatch.NewFeat(
			s.GetPort(),
			s.GetJoinTarget(),
			s.GetDispatchService(),
		)
	}
	return s.FeatDispatch
}

func require[T any](arg T) T {
	check(arg)
	return arg
}

func check(arg any) {
	if arg == nil || arg == "" || reflect.ValueOf(arg).IsZero() {
		log.Panicf("nil dependency: %T", arg)
	}
}
