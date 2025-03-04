package main

import (
	"context"
	_ "github.com/cryptopunkscc/portal/runtime/apphost"
)

func main() {
	t := Tray{}
	ctx := context.Background()
	err := t.Run(ctx)
	if err != nil {
		panic(err)
	}
}
