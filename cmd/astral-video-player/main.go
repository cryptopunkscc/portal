package main

import (
	"context"

	player "github.com/cryptopunkscc/portal/apps/player/src"
	"github.com/cryptopunkscc/portal/apps/player/vlc"
)

func main() {
	service := player.Service{
		Name: "video",
	}
	var err error
	if service.Player, err = vlc.NewPlayer(); err != nil {
		panic(err)
	}
	if err = service.Serve(context.Background()); err != nil {
		panic(err)
	}
}
