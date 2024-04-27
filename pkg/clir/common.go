package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/launcher"
	"github.com/cryptopunkscc/go-astral-js/pkg/cmd/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

func cliOpen(ctx context.Context, bindings runtime.New) func(f *FlagsPath) (err error) {
	return func(f *FlagsPath) (err error) {
		return open.Run(ctx, bindings, f.Path)
	}
}

func cliLauncher(ctx context.Context, bindings runtime.New) func() error {
	return func() error { return launcher.Run(ctx, bindings) }
}
