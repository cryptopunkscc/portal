//go:build unix

package main

import (
	"github.com/cryptopunkscc/portal/pkg/fs2"
	"os"
	"path/filepath"
)

func binariesDir() string {
	if fs2.CanWriteToDir("/bin") {
		return "/bin"
	}

	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".local/bin")
	}

	panic(err)
}
