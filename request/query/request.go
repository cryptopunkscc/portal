package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/rpc"
)

type Requester target.Port

var Request = Requester(target.PortOpen)

func (port Requester) Run(ctx context.Context, src string) (err error) {
	log := plog.Get(ctx).Type(port)
	log.Println("Running query", port, src)
	flow, err := rpc.QueryFlow(id.Anyone, port.Base)
	if err != nil {
		return
	}
	flow.Logger(log)
	defer flow.Close()
	err = rpc.Command(flow, port.Name, src)
	if err != nil {
		log.E().Printf("cannot query %s %s: %v", port, src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	c := make(chan any)
	go func() {
		_ = rpc.Await(flow)
		close(c)
	}()
	select {
	case <-ctx.Done():
	case <-c:
	}
	return
}

func (port Requester) Start(ctx context.Context, src string) (err error) {
	plog.Get(ctx).Type(port).Println("starting query", port, src)
	request := rpc.NewRequest(id.Anyone, port.Base)
	err = rpc.Command(request, port.Name, src)
	if err != nil {
		plog.Get(ctx).Type(port).E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	return
}
