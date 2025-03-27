//go:build debug || test

package plog

import (
	"os"
)

func init() {
	if _, ok := os.LookupEnv("PLOG"); !ok {
		Verbosity = all
	}
}
