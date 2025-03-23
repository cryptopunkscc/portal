//go:build windows

package main

import (
	"os"
)

func init() {
	panic("TODO")
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
}
