package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/exec"
	"github.com/cryptopunkscc/portal/runner/webview"
	"log"
	"os"
	"path"
)

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

	webview.Run(ctx, src, title)
}
