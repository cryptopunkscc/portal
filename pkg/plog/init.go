package plog

import (
	"os"
	"strconv"
)

func init() {
	if level, ok := os.LookupEnv("PLOG"); ok {
		if i, err := strconv.Atoi(level); err == nil {
			Verbosity = Level(i)
		}
	}
}
