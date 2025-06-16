package rpc

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

type Router interface {
	Init(ctx context.Context) error
	Listen() error
}

func Start(router Router) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		if err = router.Init(ctx); err != nil {
			return
		}
		go func() {
			if err := router.Listen(); err != nil {
				plog.Get(ctx).Type(router).E().Println(err)
			}
		}()
		return
	}
}

func Run(router Router) func(ctx context.Context) (err error) {
	return func(ctx context.Context) (err error) {
		if err = router.Init(ctx); err != nil {
			return
		}
		return router.Listen()
	}
}
