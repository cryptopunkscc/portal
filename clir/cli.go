package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
)

type cli struct {
	clir     *clir.Cli
	ctx      context.Context
	bindings runtime.New
}

func newCli(ctx context.Context, bindings runtime.New) *cli {
	return &cli{ctx: ctx, bindings: bindings}
}
