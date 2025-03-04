package main

import (
	"path/filepath"
	"strings"
)

func fixCmd(args []string) []string {
	for i, arg := range args {
		args[i] = fixPath(arg)
	}
	return args
}

func fixPath(str string) string {
	if strings.HasPrefix(str, "./") || strings.HasPrefix(str, "../") {
		abs, err := filepath.Abs(str)
		if err != nil {
			panic(err)
		}
		return abs
	}
	return str
}
