package dev

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/go-astral-js/feat/apps"
	"github.com/cryptopunkscc/go-astral-js/pkg/portal"
	"github.com/cryptopunkscc/go-astral-js/pkg/rpc"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/cryptopunkscc/go-astral-js/pkg/target"
	"github.com/cryptopunkscc/go-astral-js/runner/goja"
	"github.com/cryptopunkscc/go-astral-js/runner/goja_dev"
	"github.com/cryptopunkscc/go-astral-js/runner/serve"
	"github.com/cryptopunkscc/go-astral-js/runner/wails"
	"github.com/cryptopunkscc/go-astral-js/runner/wails_dev"
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
		Resolve:  resolvePortals,
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
	prefix := "dev"
	switch v := t.(type) {
	case target.App:
		typ := v.Type()
		switch {
		case typ.Is(target.Backend):
			return goja.Run(ctx, bindings, v, prefix)
		case typ.Is(target.Frontend):
			return wails.Run(bindings, v, prefix)
		default:
			return fmt.Errorf("invalid app target: %v", v.Path())
		}
	case target.Project:
		typ := v.Type()
		switch {
		case typ.Is(target.Backend):
			return goja_dev.NewBackend(ctx, bindings, v).Start()
		case typ.Is(target.Frontend):
			return wails_dev.NewFrontend(bindings, v).Start()
		default:
			return fmt.Errorf("invalid dev target: %v", t.Path())
		}
	}
	return fmt.Errorf("invalid target type %T", t)
}

var resolvePortals = portal.ResolvePortals

var Handlers = rpc.Handlers{
	"ping":    func() {},
	"open":    portal.NewCmdOpener(resolvePortals, action).Open,
	"observe": apps.Observe,
}

const action = "o"
