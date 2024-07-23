package target

import (
	"os"
	"path/filepath"
)

func Abs(src string) string {
	if filepath.IsAbs(src) {
		return src
	}
	base, err := os.Getwd()
	if err != nil {
		return src
	}
	return filepath.Join(base, src)
}
