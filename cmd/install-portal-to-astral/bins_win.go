//go:build windows

package main

import (
	"os"
	"path/filepath"
)

func binariesDir() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cache, "portal", "bin")
}
