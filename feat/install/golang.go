package install

import (
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
)

// check if go is installed
func checkGo() {
	log.Println("go version")
	if err := exec.Run(".", "go", "version"); err != nil {
		log.Fatalln("required go binary have to be installed manually")
	}
}
