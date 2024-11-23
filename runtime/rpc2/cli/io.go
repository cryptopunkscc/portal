package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type ReadWriteCloser struct {
	io.Reader
	io.Writer
	close func()
}

func (r ReadWriteCloser) Close() error {
	r.close()
	return nil
}

type Marshaler interface {
	MarshalCLI() string
}

func Marshal(a any) (result []byte, err error) {
	if a == stream.End {
		return
	}
	switch t := a.(type) {
	case Marshaler:
		result = []byte(t.MarshalCLI())
	case string:
		result = []byte(t)
	case []byte:
		result = t
	case error:
		result = fmt.Appendf(nil, "%e", t)
	case stream.Failure:
		result = fmt.Appendf(nil, "API: %s", t.Error)
	default:
		result, err = json.Marshal(a)
	}
	result = bytes.TrimSuffix(result, []byte("\n"))
	return
}
