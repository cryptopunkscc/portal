package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"log"
)

func buildPortal() {
	if err := exec.Run(
		".", "go", "install",
		"-tags", "desktop,wv2runtime.download,production",
		"-ldflags", "-w -s",
		"./cmd/portal",
	); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println()
	log.Println("portal installed successfully")
}

func buildPortalDev() {
	if err := exec.Run(
		".", "go", "install",
		"-tags", "dev",
		"./cmd/portal-dev",
	); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println()
	log.Println("portal dev installed successfully")
}
