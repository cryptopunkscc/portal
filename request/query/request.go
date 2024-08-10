package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
)

type Requester target.Port

var Request = Requester(target.PortOpen)

func (port Requester) Run(ctx context.Context, src string, args ...string) (err error) {
	log := plog.Get(ctx).Type(port)
	log.Println("Running query", port, src, args)
	flow, err := rpc.QueryFlow(id.Anyone, port.Base)
	if err != nil {
		return
	}
	flow.Logger(log)
	if err != nil {
		return err
	}
	defer flow.Close()
	typ := target.ParseType(target.TypeAny, args...)
	sTyp := fmt.Sprintf("%d", typ)
	err = rpc.Command(flow, port.Name, src, sTyp)
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

func (port Requester) Start(ctx context.Context, src string, args ...string) (err error) {
	plog.Get(ctx).Type(port).Println("starting query", port, src, args)
	request := rpc.NewRequest(id.Anyone, port.Base)
	typ := target.ParseType(target.TypeAny, args...)
	sTyp := fmt.Sprintf("%d", typ)
	err = rpc.Command(request, port.Name, src, sTyp)
	if err != nil {
		plog.Get(ctx).Type(port).E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	return
}
