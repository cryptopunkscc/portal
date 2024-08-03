//go:build !windows

package main

import (
	"os"
	"path/filepath"
)

func binariesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".local/bin")
}
