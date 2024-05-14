package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/feat/serve"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"log"
	"reflect"
)

func Run(
	ctx context.Context,
	bindings runtime.New,
	src string,
	attach bool,
) (err error) {
	return portal.Runner[target.Portal]{
		Action:   action,
		Port:     "dev.portal",
		New:      bindings,
		Serve:    serve.Run,
		Resolve:  portal.ResolvePortals,
		Attach:   Attach,
		Handlers: Handlers,
	}.Run(ctx, src, attach)
}

func Attach(
	ctx context.Context,
	bindings runtime.New,
	t target.Portal,
	_ ...string,
) (err error) {
	log.Println("dev attach", reflect.TypeOf(t), t)
	switch v := t.(type) {
	case target.App:
		return open.Attach(ctx, bindings, v, "dev")
	case target.Project:
		return AttachDev(ctx, bindings, v)
	}
	return fmt.Errorf("invalid target type %T", t)
}

func AttachDev(
	ctx context.Context,
	bindings runtime.New,
	t target.Project,
) (err error) {
	typ := t.Type()
	switch {
	case typ.Is(target.Backend):
		return NewBackend(ctx, bindings, t).Start()
	case typ.Is(target.Frontend):
		return NewFrontend(bindings, t).Start()
	default:
		return fmt.Errorf("invalid target: %v", t.Path())
	}
}

var Handlers = rpc.Handlers{
	"ping":    func() {},
	"open":    portal.NewCmdOpener(portal.ResolvePortals, action).Open,
	"observe": apps.Observe,
}

const action = "o"
