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
	m.Resolve = Any[Runnable](
		goja.Runner(m.cores().NewBackendFunc()).Try,
		m.htmlRunner().Try,
		exec2.DistRunner.Try,
		exec2.BundleRunner.Try,
	)
	return m
}
