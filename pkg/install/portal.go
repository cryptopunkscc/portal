package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"log"
)

func buildPortalDev() {
	if err := exec.Run(".", "go", "install", "-tags", "dev", "./cmd/portal"); err != nil {
		log.Fatalln("portal dev install failed: ", err)
	}
	log.Println()
	log.Println("portal dev installed successfully")
}
