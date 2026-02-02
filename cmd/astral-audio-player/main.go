package main

import (
	"context"

	"github.com/cryptopunkscc/portal/apps/player/beep"
	"github.com/cryptopunkscc/portal/apps/player/src"
)

func main() {
	service := player.Service{
		Name:   "audio",
		Player: &beep.Player{},
	}
	if err := service.Serve(context.Background()); err != nil {
		panic(err)
	}
}
