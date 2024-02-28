//go:build !dev

package main

import (
	"github.com/leaanthony/clir"
	"log"
)

func main() {
	cli := clir.NewCli(PortalName, PortalProdDescription, PortalVersion)
	flags := &FlagsPath{}
	cli.AddFlags(flags)
	cli.Action(func() error { return cliApplication(flags) })
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
