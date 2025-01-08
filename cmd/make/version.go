package main

import (
	"github.com/cryptopunkscc/portal/runner/version"
	"os"
)

func resolveVersion() {
	file, err := os.Create("./runner/version/name")
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
	file, err := os.Create("./runner/version/name")
	if err != nil {
		panic(err)
	}
	if err = file.Truncate(0); err != nil {
		panic(err)
	}
}
