package main

import (
	"context"
	"github.com/cryptopunkscc/go-astral-js/pkg/backend/goja"
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

	goja.RunSource(src)

	ctx := context.Background()
	<-ctx.Done()
}
