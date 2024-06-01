package msg

import (
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
)

func NewSend(port target.Port) func(msg target.Msg) error {
	request := rpc.NewRequest(id.Anyone, port.String())
	return func(msg target.Msg) error {
		return rpc.Command(request, "", msg)
	}
}
