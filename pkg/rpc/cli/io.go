package cli

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"io"
	"os"
	"strings"
)

func cliConnection() stream.Serializer {
	args := strings.Join(os.Args[1:], " ")
	plog.D().Println(args)
	return stream.Serializer{
		Reader: io.MultiReader(
			strings.NewReader(args+"\n"),
			os.Stdin,
		),
		Writer: os.Stdout,
		Closer: stream.Closer(func() error {
			os.Exit(0)
			return nil
		}),
		Codec: stream.Codec{
			Marshal: Marshal,
		},
	}
}
