//go:build linux && !android

package dir

import (
	"os"
)

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir = mk(dir, ".local/share/portald")
	Init(dir)
}
