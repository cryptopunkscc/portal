//go:build !dev

package clir

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/runner"
	"github.com/leaanthony/clir"
	"log"
)

func Run(bindings runner.Bindings) {
	cli := clir.NewCli(PortalName, PortalProdDescription, PortalVersion)
	flags := &FlagsPath{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cliApplication(bindings)(flags) })
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
