package target

import (
	"os"
	"path"
)

func Abs(src string) string {
	if path.IsAbs(src) {
		return src
	}
	base, err := os.Getwd()
	if err != nil {
		return src
	}
	return path.Join(base, src)
}
