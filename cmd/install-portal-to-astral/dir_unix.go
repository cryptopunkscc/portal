//go:build unix

package main

import (
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/pkg/fs2"
)

func Dir() string {
	if fs2.CanWriteToDir("/bin") {
		return "/bin"
	}

	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".local/bin")
	}

	panic(err)
}
