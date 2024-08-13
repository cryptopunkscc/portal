package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/msg"
	"github.com/cryptopunkscc/portal/target"
)

func Mutable[T target.Portal_](
	newRuntime target.NewRuntime,
	portMsg target.Port,
	newRunner func(target.NewRuntime, target.MsgSend) target.Runner[T],
) target.Run[target.Portal_] {
	return runner[T]{
		portMsg:    portMsg,
		newRuntime: newRuntime,
		newRunner:  newRunner,
	}.Run
}

func Immutable[T target.Portal_](
	newRuntime target.NewRuntime,
	portMsg target.Port,
	newRunner func(target.NewRuntime) target.Runner[T],
) target.Run[target.Portal_] {
	return runner[T]{
		portMsg:    portMsg,
		newRuntime: newRuntime,
		newRunner: func(api target.NewRuntime, _ target.MsgSend) target.Runner[T] {
			return newRunner(api)
		},
	}.Run
}

type runner[T target.Portal_] struct {
	portMsg    target.Port
	newRuntime target.NewRuntime
	newRunner  func(target.NewRuntime, target.MsgSend) target.Runner[T]
}

func (r runner[T]) Run(ctx context.Context, portal target.Portal_) (err error) {
	t, ok := portal.(T)
	if !ok {
		return target.ErrNotTarget
	}

	var reloader msg.Reloader
	client := msg.NewClient(r.portMsg)
	sendMsg := client.Send
	newRuntime := func(ctx context.Context, portal target.Portal_) target.Runtime {
		api := r.newRuntime(ctx, portal)
		if api != nil {
			client.Init(reloader, api)
		}
		if err = client.Connect(ctx, t); err != nil {
			plog.Get(ctx).Type(r).P().Println(err)
		}
		return api
	}
	_runner := r.newRunner(newRuntime, sendMsg)
	reloader = _runner
	return _runner.Run(ctx, t)
}
