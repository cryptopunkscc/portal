package main

import (
	"github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
)

func installAstral() {
	if err := exec.Run(".", "which", "astrald"); err == nil {
		_ = exec.Run(".", "astrald", "-v")
		return
	}
	astrald := "github.com/cryptopunkscc/astrald/cmd/astrald@" + portal.AstralVersion
	log.Println("Installing", astrald)
	if err := exec.Run(".", "go", "install", astrald); err != nil {
		log.Fatalln("cannot install astrald:", err)
	}
	log.Println("Astral installed")
}
