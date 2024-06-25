package main

import (
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
)

// check if npm is installed
func checkNpm() {
	log.Println("npn -v")
	if err := exec.Run(".", "npm", "-v"); err != nil {
		log.Fatalln("required go binary have to be installed manually")
	}
}
