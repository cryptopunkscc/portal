package main

import (
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/factory/bind"
	"github.com/cryptopunkscc/portal/resolve/sources"
	"github.com/cryptopunkscc/portal/runner/cli"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/go_dev"
	"github.com/cryptopunkscc/portal/runner/multi"
	"github.com/cryptopunkscc/portal/runner/open"
	"github.com/cryptopunkscc/portal/runner/reload"
	"github.com/cryptopunkscc/portal/runner/version"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
	"github.com/cryptopunkscc/portal/runtime/rpc/cmd"
)

func main() { cli.Run(Application[ProjectGo]{}.handler()) }

type Application[T ProjectGo] struct{}

func (a Application[T]) handler() cmd.Handler {
	return cmd.Handler{
		Func: open.Runner[T](a),
		Name: "dev-go",
		Desc: "Start portal golang app development.",
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
		reload.Mutable(bind.DefaultRuntime(), go_dev.Adapter(exec.DistRun)),
	)
}

func (a Application[T]) Resolver() Resolve[T] { return sources.Resolver[T]() }
