//go:build windows

package dir

import (
	"os"
)

func portaldDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	Init(dir)
}
