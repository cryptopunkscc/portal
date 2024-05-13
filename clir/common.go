package clir

import (
	"context"
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/feat/launcher"
	"github.com/cryptopunkscc/go-astral-js/feat/open"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
	"log"
)

type FlagsPath struct {
	Path string `pos:"1" default:"."`
}

type FlagsOpen struct {
	FlagsPath
	Attach bool `name:"attach" description:"Attach runner to the current process instead of dispatching to execution portal service. If given path contains multiple portals each will be run as a child process."`
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

func (c Cli) Version(_ *struct{}) (_ error) {
	log.Println(portal.GoModuleVersion())
	return
}
