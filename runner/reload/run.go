package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/pkg/plog"
)

func Mutable[T target.Portal_](
	newCore bind.NewCore,
	newReRunner func(bind.NewCore, target.MsgSend) target.ReRunner[T],
) target.Run[target.Portal_] {
	return runner(newCore, newReRunner)
}

func Immutable[T target.Portal_](
	newCore bind.NewCore,
	newReRunner func(bind.NewCore) target.ReRunner[T],
) target.Run[target.Portal_] {
	return runner(newCore, func(api bind.NewCore, _ target.MsgSend) target.ReRunner[T] {
		return newReRunner(api)
	})
}

func runner[T target.Portal_](
	newCore bind.NewCore,
	newReRunner func(bind.NewCore, target.MsgSend) target.ReRunner[T],
) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		t, ok := src.(T)
		if !ok {
			return target.ErrNotTarget
		}

		var reload Reload
		client := newClient()
		sendMsg := client.Send
		newCore := func(ctx context.Context, portal target.Portal_) (bind.Core, context.Context) {
			core, ctx := newCore(ctx, portal)
			if core != nil {
				client.Init(reload, core)
			}
			if err = client.Connect(ctx, t); err != nil {
				plog.Get(ctx).Scope("ReRunner").E().Println(err)
			}
			return core, ctx
		}
		_runner := newReRunner(newCore, sendMsg)
		reload = _runner.Reload
		return _runner.Run(ctx, t, args...)
	}
}
