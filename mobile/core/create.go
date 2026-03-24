package core

import (
	"os"
	"path"

	"github.com/cryptopunkscc/astrald/astral"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	ip "github.com/cryptopunkscc/astrald/mod/ip/src"
	"github.com/cryptopunkscc/portal/mobile"
	"github.com/cryptopunkscc/portal/pkg/util/plog"
)

func Create(api mobile.Api) mobile.Core {
	plog.Verbosity = 100
	ip.InterfaceAddrs = interfaceAddrsFunc(api)
	ether.NetInterfaces = netInterfacesFunc(api)

	err := os.Chdir(api.DataDir())
	if err != nil {
		panic(err)
	}

	m := &Service{}
	m.api = api
	m.ctx, m.cancel = astral.NewContext(nil).WithCancel()
	//m.astrald.Ghost = true
	m.astrald.DBRoot = api.DataDir()
	m.astrald.NodeRoot = path.Join(api.DataDir(), "astrald")

	return m
}
