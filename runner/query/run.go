package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
)

type Runner[T target.Portal] struct {
	port string
}

func NewRunner[T target.Portal](port string) *Runner[T] {
	return &Runner[T]{port: port}
}

func (r Runner[T]) Run(ctx context.Context, src string, args ...string) (err error) {
	flow, err := rpc.QueryFlow(id.Anyone, r.port)
	if err != nil {
		return err
	}
	defer flow.Close()
	typ := target.ParseType(target.TypeAny, args...)
	sTyp := fmt.Sprintf("%d", typ)
	err = rpc.Command(flow, "", src, sTyp)
	if err != nil {
		plog.Get(ctx).Type(r).E().Printf("cannot query %s: %v", src, err)
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

func (r Runner[T]) Start(ctx context.Context, src string, args ...string) (err error) {
	plog.Get(ctx).Type(r).Println("starting query", src, args)
	request := rpc.NewRequest(id.Anyone, r.port)
	typ := target.ParseType(target.TypeAny, args...)
	sTyp := fmt.Sprintf("%d", typ)
	err = rpc.Command(request, "", src, sTyp)
	if err != nil {
		plog.Get(ctx).Type(r).E().Printf("cannot query %s: %v", src, err)
		return fmt.Errorf("cannot query %s: %w", src, err)
	}
	return
}
