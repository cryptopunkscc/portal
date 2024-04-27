//go:build !dev

package clir

import (
	"github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/runtime"
	"github.com/leaanthony/clir"
	"log"
)

func Run(ctx context.Context, bindings runtime.New) {
	cli := clir.NewCli(portal.Name, portal.ProdDescription, portal.Version)
	flags := &FlagsPath{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cliApplication(cts, bindings)(flags) })
	cli.NewSubCommand("launcher", "Start portal launcher GUI.").Action(cliLauncher(ctx, bindings))
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
