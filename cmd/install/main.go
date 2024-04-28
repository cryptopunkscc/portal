package main

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/install"
	"log"
	"os"
	"strconv"
)

func main() {
	deps := install.All
	if len(os.Args) > 1 {
		arg, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		deps = install.Dependency(arg)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dev := install.NewPortalDev(wd)
	dev.Install(deps)
}
