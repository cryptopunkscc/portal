package cli

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/rpc/stream"
	"gopkg.in/yaml.v3"
)

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
		result, err = yaml.Marshal(a)
	}
	if !bytes.HasSuffix(result, []byte("\n")) {
		result = append(result, '\n')
	}
	return
}
