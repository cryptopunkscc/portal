package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc"
	"github.com/cryptopunkscc/portal/target"
)

type Open struct {
	port target.Port
}

func NewOpen() *Open {
	return &Open{port: target.PortOpen}
}

func (r Open) Run(ctx context.Context, src string, args ...string) (err error) {
	log := plog.Get(ctx).Type(r)
	log.Println("Running query", r.port, src, args)
	flow, err := rpc.QueryFlow(id.Anyone, r.port.Base)
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
	err = rpc.Command(flow, r.port.Name, src, sTyp)
	if err != nil {
		log.E().Printf("cannot query %s %s: %v", r.port, src, err)
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

func (r Open) Start(ctx context.Context, src string, args ...string) (err error) {
	plog.Get(ctx).Type(r).Println("starting query", r.port, src, args)
	request := rpc.NewRequest(id.Anyone, r.port.Base)
	typ := target.ParseType(target.TypeAny, args...)
	sTyp := fmt.Sprintf("%d", typ)
	err = rpc.Command(request, r.port.Name, src, sTyp)
	if err != nil {
		plog.Get(ctx).Type(r).E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	return
}
