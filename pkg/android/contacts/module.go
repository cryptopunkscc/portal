package contacts

import (
	"github.com/cryptopunkscc/astrald/log"
	"github.com/cryptopunkscc/astrald/node/assets"
	"github.com/cryptopunkscc/astrald/node/modules"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
)

const ServiceName = "contacts.*"

func init() {
	if err := modules.RegisterModule(ServiceName, Loader{}); err != nil {
		panic(err)
	}
}

type Loader struct{}

func (Loader) Load(node modules.Node, _ assets.Assets, log *log.Logger) (modules.Module, error) {
	module := rpc.NewModule(node, ServiceName)
	module.Interface(&service{node: node, log: log})
	return module, nil
}
