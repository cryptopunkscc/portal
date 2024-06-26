package main

import (
	"github.com/cryptopunkscc/portal/feat/version"
	"os"
)

func resolveVersion() {
	file, err := os.Create("./feat/version/name")
	if err != nil {
		return
	}
	defer file.Close()
	name := version.Resolve()
	if _, err = file.WriteString(name); err != nil {
		panic(err)
	}
}

func clearVersion() {
	file, err := os.Create("./feat/version/name")
	if err != nil {
		panic(err)
	}
	if err = file.Truncate(0); err != nil {
		panic(err)
	}
}
