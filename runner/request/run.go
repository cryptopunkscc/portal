package request

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc2"
	apphostRpc "github.com/cryptopunkscc/portal/runtime/rpc2/apphost"
)

type Requester struct{ apphost.Port }

var Open = Requester{target.PortOpen}

func (port Requester) Start(ctx context.Context, src string, _ ...string) (err error) {
	log := plog.Get(ctx).Type(port)
	log.Println("starting query", port, src)
	request := apphostRpc.Default().Request("portal", port.Base())
	err = rpc.Command(request, port.Name(), src)
	if err != nil {
		log.E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	log.Println("started query", port, src)
	return
}
