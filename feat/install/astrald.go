package install

import (
	portal "github.com/cryptopunkscc/go-astral-js"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
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
