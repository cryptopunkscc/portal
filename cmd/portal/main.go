package main

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/feat/start"
	"github.com/cryptopunkscc/portal/feat/version"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
	"github.com/cryptopunkscc/portal/request/exec"
	"github.com/cryptopunkscc/portal/request/query"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cli"
	"github.com/cryptopunkscc/portal/runtime/rpc2/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := plog.New().D().Scope("app").Set(&ctx)
	go sig.OnShutdown(cancel)

	err := cli.New(cmd.Handler{
		Name: "portal",
		Desc: "Portal command line.",
		Func: start.Feat(deps{}),
		Params: cmd.Params{
			{Type: "string", Desc: "Application source. The source can be a app name, package name, app bundle path or app dir."},
		},
		Sub: cmd.Handlers{
			{Name: "v", Desc: "Print version", Func: version.Run},
		},
	}).Run(ctx)

	if err != nil {
		log.Println(err)
	}
	cancel()
}

type deps struct{}

func (m deps) Port() apphost.Port      { return target.PortPortal }
func (m deps) Serve() target.Request   { return exec.Request("portal-app") }
func (m deps) Request() target.Request { return query.Request.Run }
