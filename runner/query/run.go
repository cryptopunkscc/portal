package query

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/go-astral-js/pkg/plog"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/target"
	"strings"
)

const port = "portal.open"

type Runner[T target.Portal] struct {
	prefix []string
}

func NewRunner[T target.Portal](prefix ...string) *Runner[T] {
	return &Runner[T]{prefix: prefix}
}

func (r Runner[T]) Run(ctx context.Context, src string, _ ...string) (err error) {
	srv := strings.Join(append(r.prefix, port), ".")
	flow, err := rpc.QueryFlow(id.Anyone, srv)
	if err != nil {
		return err
	}
	defer flow.Close()
	err = rpc.Command(flow, "", src)
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
