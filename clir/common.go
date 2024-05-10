package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/feat/launcher"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

type FlagsOpen struct {
	FlagsPath
	Attach bool `name:"attach" description:"Attach execution to the current process instead of dispatching to portal service."`
}

type Cli struct {
	*clir.Cli
	ctx      context.Context
	bindings runtime.New
}

func (c Cli) Open(f *FlagsOpen) (err error) {
	return open.Run(c.ctx, c.bindings, f.Path, f.Attach)
}

func (c Cli) Launcher() error {
	return launcher.Run(c.ctx, c.bindings)
}
