package main

import (
	"github.com/cryptopunkscc/portal"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	"os"
	"path/filepath"
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

func buildAstral() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("cannot resolve working dir:", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("cannot resolve home dir:", err)
	}
	astrald := "github.com/cryptopunkscc/astrald@" + portal.AstralVersion
	astrald = filepath.Join(home, "go/pkg/mod", astrald)
	out := filepath.Join(wd, "cmd/portal-installer/bin/")
	if err := exec.Run(
		astrald, "go", "build",
		"-o", out,
		"./cmd/astrald",
	); err != nil {
		log.Fatalln("cannot build astrald:", err)
	}
	log.Println("Astral build succeeded")
}
