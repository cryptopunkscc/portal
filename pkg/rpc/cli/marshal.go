package cli

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/cryptopunkscc/portal/pkg/plog"
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
		var es plog.ErrStack
		if ok := errors.As(t.Error, &es); ok {
			result = fmt.Appendf(nil, "ERROR: %s\n\n%s", es.Error(), es.Stack())
		} else {
			result = fmt.Appendf(nil, "ERROR: %s", t.Error)
		}
	default:
		result, err = yaml.Marshal(a)
	}
	if !bytes.HasSuffix(result, []byte("\n")) {
		result = append(result, '\n')
	}
	return
}
