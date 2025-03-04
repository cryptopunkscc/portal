package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/portal/runtime/rpc2/stream"
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
		result = fmt.Appendf(nil, "%e\n", t)
	case stream.Failure:
		result = fmt.Appendf(nil, "API: %s\n", t.Error)
	default:
		result, err = json.Marshal(a)
	}
	return
}
