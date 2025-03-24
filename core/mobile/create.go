package core

import (
	"context"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/api/mobile"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/html"
	"github.com/cryptopunkscc/portal/resolve/js"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/goja"
)

func Create(api mobile.Api) mobile.Core {
	env.AstraldHome.SetDir(api.DataDir(), "astrald")
	env.AstraldDb.SetDir(api.DbDir())
	env.PortaldApps.SetDir(api.DataDir(), "portald", "apps")
	env.PortaldTokens.SetDir(api.DataDir(), "portald", "tokens")

	tcp.InterfaceAddrs = interfaceAddrsFunc(api)
	ether.NetInterfaces = netInterfacesFunc(api)
	plog.Verbosity = 100

	m := &service{}
	ctx := context.Background()
	astraldHomeDir := mem.NewVar(env.AstraldHome.MkdirAll())
	astraldDbDir := mem.NewVar(env.AstraldDb.MkdirAll())
	tokensDir := mem.NewVar(env.PortaldTokens.MkdirAll())
	appsDir := mem.NewVar(env.PortaldApps.MkdirAll())

	m.ctx = ctx
	m.mobile = api
	m.TokensDir = tokensDir
	m.NodeDir = astraldHomeDir
	m.AppsDir = appsDir
	m.Astrald = &astrald{
		NodeRoot: astraldHomeDir,
		DbRoot:   astraldDbDir,
	}
	m.CreateTokens = []string{
		"portal.launcher",
	}
	m.Resolve = Any[App_](
		Skip("node_modules"),
		Try(js.ResolveDist),
		Try(html.ResolveDist),
		Try(exec.ResolveDist),
	)
	m.Runners = func(schemaPrefix []string) []Run[Portal_] {
		return []Run[Portal_]{
			app.Runner(goja.NewRun(m.cores().NewBackendFunc())),
			app.Runner(m.htmlRun),
		}
	}
	m.find = FindByPath(source.File, m.Resolve).
		OrById(path.Resolver(m.Resolve, env.PortaldApps.Source()))
	return m
}
