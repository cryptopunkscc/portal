//go:build legacy

package webview

//package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/sig"
)

// legacy main function
func main() {
	file := os.Args[1]
	srcBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	src := string(srcBytes)
	title := filepath.Base(file)

	ctx, cancel := context.WithCancel(context.Background())
	go sig.OnShutdown(plog.New(), cancel)

	Run(ctx, src, title)
}
