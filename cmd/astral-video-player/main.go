package main

import (
	"context"

	"github.com/cryptopunkscc/portal/apps/player/video"
	"github.com/cryptopunkscc/portal/apps/player/vlc"
)

func main() {
	service := video.Service{}
	var err error
	if service.Player, err = vlc.NewPlayer(); err != nil {
		panic(err)
	}
	if err = service.Serve(context.Background()); err != nil {
		panic(err)
	}
}
