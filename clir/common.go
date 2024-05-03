package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/launcher"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

type FlagsOpen struct {
	FlagsPath
	Attach bool `name:"attach" description:"Attach execution to the current process instead of dispatching to portal service."`
}

func cliOpen(ctx context.Context, bindings runtime.New) func(f *FlagsOpen) (err error) {
	return func(f *FlagsOpen) (err error) {
		return open.Run(ctx, bindings, f.Path, f.Attach)
	}
}

func cliLauncher(ctx context.Context, bindings runtime.New) func() error {
	return func() error { return launcher.Run(ctx, bindings) }
}
