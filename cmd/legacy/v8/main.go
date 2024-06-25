package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/v8"
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

	// prepare context
	ctx, cancel := context.WithCancel(context.Background())
	exec.OnShutdown(cancel)

	if err = v8.Run(ctx, file, src); err != nil {
		log.Fatalln(err)
	}
	<-ctx.Done()
}
