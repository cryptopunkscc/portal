package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/bind"
)

func Mutable[T target.Portal_](
	newRuntime bind.NewRuntime,
	newReRunner func(bind.NewRuntime, target.MsgSend) target.ReRunner[T],
) target.Run[target.Portal_] {
	return runner(newRuntime, newReRunner)
}

func Immutable[T target.Portal_](
	newRuntime bind.NewRuntime,
	newReRunner func(bind.NewRuntime) target.ReRunner[T],
) target.Run[target.Portal_] {
	return runner(newRuntime, func(api bind.NewRuntime, _ target.MsgSend) target.ReRunner[T] {
		return newReRunner(api)
	})
}

func runner[T target.Portal_](
	newRuntime bind.NewRuntime,
	newReRunner func(bind.NewRuntime, target.MsgSend) target.ReRunner[T],
) target.Run[target.Portal_] {
	return func(ctx context.Context, src target.Portal_, args ...string) (err error) {
		t, ok := src.(T)
		if !ok {
			return target.ErrNotTarget
		}

		var reRun ReRun
		client := newClient()
		sendMsg := client.Send
		newRuntime := func(ctx context.Context, portal target.Portal_) bind.Runtime {
			runtime := newRuntime(ctx, portal)
			if runtime != nil {
				client.Init(reRun, runtime)
			}
			if err = client.Connect(ctx, t); err != nil {
				plog.Get(ctx).Scope("ReRunner").P().Println(err)
			}
			return runtime
		}
		_runner := newReRunner(newRuntime, sendMsg)
		reRun = _runner.ReRun
		return _runner.Run(ctx, t, args...)
	}
}
