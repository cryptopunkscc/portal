package core

import (
	"context"

	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/api/mobile"
	. "github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	exec2 "github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/runner/goja"
)

func Create(api mobile.Api) mobile.Core {
	tcp.InterfaceAddrs = interfaceAddrsFunc(api)
	ether.NetInterfaces = netInterfacesFunc(api)
	plog.Verbosity = 100

	m := &service{}
	m.mobile = api
	m.ctx = context.Background()
	m.Config.Dir = api.DataDir()
	m.Config.Astrald = "astrald"
	m.Config.Tokens = "portald/tokens"
	m.Config.Apps = "portald/apps"
	m.Config.Bin = "portald/bin"
	m.Config.Config.Node.Log.DisableColors = true
	m.Config.AstralDB = api.DbDir()
	_ = m.Configure()
	m.Astrald = &astrald{
		NodeRoot: m.Config.Astrald,
		DbRoot:   api.DbDir(),
	}
	m.ExtraTokens = []string{
		"portal.launcher",
	}
	m.Resolve = Any[Runnable](
		Skip("node_modules"),
		goja.Runner(m.cores().NewBackendFunc()).Try,
		m.htmlRunner().Try,
		exec2.Runner{Config: m.Config}.Dist().Try,
		exec2.Runner{Config: m.Config}.Bundle().Try,
	)
	return m
}
