package main

import (
	"context"

	"github.com/cryptopunkscc/portal/apps/player/audio"
	"github.com/cryptopunkscc/portal/apps/player/beep"
)

func main() {
	service := audio.Service{
		Player: &beep.Player{},
	}
	if err := service.Serve(context.Background()); err != nil {
		panic(err)
	}
}
