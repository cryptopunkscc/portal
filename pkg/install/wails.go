package install

import (
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"log"
	"os"
	"path"
)

// install wails
func installWails() {
	if err := exec.Run(".", "which", "wails"); err == nil {
		if err = exec.Run(".", "go", "install", "github.com/wailsapp/wails/v2/cmd/wails@latest"); err != nil {
			log.Fatalln("cannot install wails:", err)
		}
		// run wails doctor for installing required dependencies
		if err = exec.Run(".", "wails", "doctor"); err != nil {
			log.Fatalln(err)
		}
		err := os.Remove(path.Join(os.Getenv("GOPATH"), "bin", "wails"))
		if err != nil {
			log.Println("cannot remove wails:", err)
		}
	}
}
