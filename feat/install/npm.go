package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"log"
)

// check if npm is installed
func checkNpm() {
	log.Println("npn -v")
	if err := exec.Run(".", "npm", "-v"); err != nil {
		log.Fatalln("required go binary have to be installed manually")
	}
}
