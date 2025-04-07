//go:build debug || test

package plog

import (
	"os"
	"strconv"
)

func init() {
	if _, ok := os.LookupEnv("PLOG"); !ok {
		Verbosity = all
		_ = os.Setenv("PLOG", strconv.Itoa(int(Verbosity)))
	}
}
