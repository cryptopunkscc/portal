package webview

//package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"log"
	"os"
	"path"
)

// legacy main function
func main() {
	file := os.Args[1]
	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	src := string(srcBytes)
	title := path.Base(file)

	ctx, cancel := context.WithCancel(context.Background())
	go exec.OnShutdown(cancel)

	Run(ctx, src, title)
}
