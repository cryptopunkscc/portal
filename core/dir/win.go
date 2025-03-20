//go:build windows

package dir

import (
	"os"
)

func portaldDir() string {
	panic("FIXME!")
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	Init(dir)
}
