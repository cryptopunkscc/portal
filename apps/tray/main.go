package main

import (
	"context"
)

func main() {
	t := Tray{}
	ctx := context.Background()
	err := t.Run(ctx)
	if err != nil {
		panic(err)
	}
}
