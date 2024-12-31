package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/bind"
	"github.com/cryptopunkscc/portal/runtime/msg"
)

func Mutable[T target.Portal_](
	newRuntime bind.NewRuntime,
	portMsg apphost.Port,
	newRunner func(bind.NewRuntime, target.MsgSend) target.Runner[T],
) target.Run[target.Portal_] {
	return runner[T]{
		portMsg:    portMsg,
		newRuntime: newRuntime,
		newRunner:  newRunner,
	}.Run
}

func Immutable[T target.Portal_](
	newRuntime bind.NewRuntime,
	portMsg apphost.Port,
	newRunner func(bind.NewRuntime) target.Runner[T],
) target.Run[target.Portal_] {
	return runner[T]{
		portMsg:    portMsg,
		newRuntime: newRuntime,
		newRunner: func(api bind.NewRuntime, _ target.MsgSend) target.Runner[T] {
			return newRunner(api)
		},
	}.Run
}

type runner[T target.Portal_] struct {
	portMsg    apphost.Port
	newRuntime bind.NewRuntime
	newRunner  func(bind.NewRuntime, target.MsgSend) target.Runner[T]
}

func (r runner[T]) Run(ctx context.Context, portal target.Portal_) (err error) {
	t, ok := portal.(T)
	if !ok {
		return target.ErrNotTarget
	}

	var reloader msg.Reloader
	client := msg.NewClient(r.portMsg)
	sendMsg := client.Send
	newRuntime := func(ctx context.Context, portal target.Portal_) bind.Runtime {
		runtime := r.newRuntime(ctx, portal)
		if runtime != nil {
			client.Init(reloader, runtime)
		}
		if err = client.Connect(ctx, t); err != nil {
			plog.Get(ctx).Type(r).P().Println(err)
		}
		return runtime
	}
	_runner := r.newRunner(newRuntime, sendMsg)
	reloader = _runner
	return _runner.Run(ctx, t)
}
