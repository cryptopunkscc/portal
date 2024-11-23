package cli

import (
	"encoding/json"
	rpc "github.com/cryptopunkscc/portal/runtime/rpc2"
	"io"
	"os"
	"strings"
)

func cliConnection() rpc.Serializer {
	args := strings.Join(os.Args[1:], " ")
	return rpc.Serializer{
		Reader: io.MultiReader(
			strings.NewReader(args+"\n"),
			os.Stdin,
		),
		Writer: os.Stdout,
		Closer: rpc.Closer(func() error {
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

func Marshal(a any) ([]byte, error) {
	switch t := a.(type) {
	case Marshaler:
		return []byte(t.MarshalCLI()), nil
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		return json.Marshal(a)
	}
}
