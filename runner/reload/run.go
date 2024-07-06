package reload

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/target"
	"github.com/cryptopunkscc/portal/target/msg"
)

func Mutable[T target.Portal](
	newApi target.NewApi,
	portMsg target.Port,
	newRunner func(target.NewApi, target.MsgSend) target.Runner[T],
) target.Run[target.Portal] {
	return Runner[T]{
		portMsg:   portMsg,
		newApi:    newApi,
		newRunner: newRunner,
	}.Run
}

func Immutable[T target.Portal](
	newApi target.NewApi,
	portMsg target.Port,
	newRunner func(target.NewApi) target.Runner[T],
) target.Run[target.Portal] {
	return Runner[T]{
		portMsg: portMsg,
		newApi:  newApi,
		newRunner: func(api target.NewApi, _ target.MsgSend) target.Runner[T] {
			return newRunner(api)
		},
	}.Run
}

type Runner[T target.Portal] struct {
	portMsg   target.Port
	newApi    target.NewApi
	newRunner func(target.NewApi, target.MsgSend) target.Runner[T]
}

func (r Runner[T]) Run(ctx context.Context, portal target.Portal) (err error) {
	t, ok := portal.(T)
	if !ok {
		return target.ErrNotTarget
	}

	var reloader msg.Reloader
	client := msg.NewClient(r.portMsg)
	sendMsg := client.Send
	newApi := func(ctx context.Context, portal target.Portal) target.Api {
		api := r.newApi(ctx, portal)
		if api != nil {
			client.Init(reloader, api)
		}
		if err = client.Connect(ctx, t); err != nil {
			plog.Get(ctx).Type(r).P().Println(err)
		}
		return api
	}
	runner := r.newRunner(newApi, sendMsg)
	reloader = runner
	return runner.Run(ctx, t)
}
