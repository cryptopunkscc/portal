//go:build !dev

package clir

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
	"log"
)

func Run(ctx context.Context, bindings runtime.New) {
	cli := Cli{
		Cli:      clir.NewCli(portal.Name, portal.ProdDescription, portal.Version),
		ctx:      ctx,
		bindings: bindings,
	}

	flags := &FlagsOpen{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cli.Open(flags) })
	cli.NewSubCommand("launcher", "Start portal launcher GUI.").Action(cli.Launcher)
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
