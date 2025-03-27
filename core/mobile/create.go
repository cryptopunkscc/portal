package core

import (
	"context"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/api/mobile"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/resolve/exec"
	"github.com/cryptopunkscc/portal/resolve/html"
	"github.com/cryptopunkscc/portal/resolve/js"
	"github.com/cryptopunkscc/portal/resolve/path"
	"github.com/cryptopunkscc/portal/resolve/source"
	"github.com/cryptopunkscc/portal/runner/app"
	"github.com/cryptopunkscc/portal/runner/goja"
	"path/filepath"
)

func Create(api mobile.Api) mobile.Core {
	tcp.InterfaceAddrs = interfaceAddrsFunc(api)
	ether.NetInterfaces = netInterfacesFunc(api)
	plog.Verbosity = 100

	m := &service{}
	m.mobile = api
	m.ctx = context.Background()
	m.Config.Astrald = filepath.Join(api.DataDir(), "astrald")
	m.Config.Tokens = filepath.Join(api.DataDir(), "portald", "tokens")
	m.Config.Apps = filepath.Join(api.DataDir(), "portald", "apps")
	m.Config.AstralDB = api.DbDir()
	m.Astrald = &astrald{
		NodeRoot: m.Config.Astrald,
		DbRoot:   api.DbDir(),
	}
	m.ExtraTokens = []string{
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
