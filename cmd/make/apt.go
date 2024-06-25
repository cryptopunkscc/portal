package main

import (
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	exec2 "os/exec"
)

// install linux apt dependencies
func installApt() {
	var deps []string

	// tray icon dependencies for linux
	deps = append(deps, trayDeps...)

	var missing []string
	for _, d := range deps {
		if err := exec2.Command("dpkg-query", "-l", d).Run(); err != nil {
			log.Printf("missing dep: %s", d)
			missing = append(missing, d)
		}
	}
	if len(missing) > 0 {
		if err := exec.Run(".", "sudo", "apt-get", "install", "gcc", "libgtk-3-dev", "libayatana-appindicator3-dev"); err != nil {
			log.Fatalln(err)
		}
	}
}
