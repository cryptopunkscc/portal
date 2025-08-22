//go:build windows

package bin

import (
	"os"
	"path/filepath"
)

func Dir() string {
	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cache, "portal", "bin")
}
