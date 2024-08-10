package target

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Request func(ctx context.Context, src string, args ...string) (err error)

func (r Request) Start(ctx context.Context, src string, args ...string) (err error) {
	go func() {
		if err = r(ctx, src, args...); err != nil {
			plog.Get(ctx).Println("Start:", err)
		}
	}()
	return nil
}
