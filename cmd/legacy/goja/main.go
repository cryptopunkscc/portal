package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/exec"
	"github.com/cryptopunkscc/go-astral-js/pkg/goja"
	"log"
	"os"
)

func main() {
	file := os.Args[1]

	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	src := string(srcBytes)

	ctx, cancel := context.WithCancel(context.Background())
	go exec.OnShutdown(cancel)

	if err = goja.NewBackend(ctx).RunSource(src); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
