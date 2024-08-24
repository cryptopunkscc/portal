package apps

import (
	"os"
	"path/filepath"
)

func DefaultDir() string {
	base, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(base, "portal", "apps")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	return dir
}
