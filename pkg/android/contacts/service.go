package contacts

import (
	"github.com/cryptopunkscc/astrald/log"
	"github.com/cryptopunkscc/astrald/node/modules"
)

type service struct {
	log  *log.Logger
	node modules.Node
}
