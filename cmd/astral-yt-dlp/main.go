package main

import (
	"context"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/src"
	"github.com/cryptopunkscc/portal/pkg/rpc/cli"
	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
)

func main() {
	cli.Run(cmd.Handler{
		Name: "astral-yt-dlp",
		Desc: "Astral youtube-dl wrapper.",
		Params: cmd.Params{
			{
				Type: "string",
				Desc: "Optional media directory. Default is current dir.",
			},
		},
		Func: func(ctx context.Context, dir string) error {
			service := astral_yt_dlp.Service{Dir: dir}
			return service.Serve(astral.NewContext(ctx))
		},
	})
}
