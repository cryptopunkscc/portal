package cli

import (
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
	"io"
	"os"
	"strings"
)

func cliConnection() stream.Serializer {
	args := strings.Join(os.Args[1:], " ")
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
		Marshal: Marshal,
	}
}
