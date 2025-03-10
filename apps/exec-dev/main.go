package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func main() { cli.Run(Application[AppExec]{}.handler()) }

type Application[T AppExec] struct{}

func (a Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[T](&a),
		Name: "dev-exec",
		Desc: "Portal development runner for executables.",
		Params: cmd.Params{
			{Type: "string", Desc: "Absolute path to app bundle or directory."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version.", Func: version.Run},
		},
	}
}

func (a Application[T]) Runner() Run[T] {
	return multi.Runner[T](
		reload.Immutable(bind.DefaultRuntime(), reload.Adapter(exec.BundleRunner(CacheDir("portal-dev")).ReRunner())),
		reload.Immutable(bind.DefaultRuntime(), reload.Adapter(exec.DistRun.ReRunner())),
	)
}

func (a Application[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }
