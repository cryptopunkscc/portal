//go:build !dev

package main

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/frontend/wails"
	"github.com/leaanthony/clir"
	"log"
)

func main() {
	cli := clir.NewCli(PortalName, PortalProdDescription, PortalVersion)
	flags := &FlagsApp{}
	cli.AddFlags(flags)
	cli.Action(func() error { return wails.Run(flags.Path, AppOptions()) })
	if err := cli.Run(); err != nil {
		log.Fatalln(err)
	}
}
