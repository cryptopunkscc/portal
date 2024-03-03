package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/runner/backend/goja"
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

	if err = goja.NewBackend().RunSource(src); err != nil {
		panic(err)
	}

	ctx := context.Background()
	<-ctx.Done()
}
