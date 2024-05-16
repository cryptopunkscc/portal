package project

import (
	"os"
	"path"
	"strings"
)

func Path(src string) (base, sub string, err error) {
	src = path.Clean(src)
	base = src
	sub = "."
	if !path.IsAbs(base) {
		sub = src
		base, err = os.Getwd()
		if err != nil {
			return
		}
	} else if strings.HasSuffix(src, ".portal") {
		base, sub = path.Split(src)
	}
	return
}

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
