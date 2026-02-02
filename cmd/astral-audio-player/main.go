package main

import (
	"context"

	"github.com/cryptopunkscc/portal/cmd/astral-audio-player/src"
)

func main() {
	if err := astral_audio_player.Serve(context.Background()); err != nil {
		panic(err)
	}
}
