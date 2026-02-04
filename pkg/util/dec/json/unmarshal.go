package json

import (
	"encoding/json"

	"github.com/cryptopunkscc/portal/pkg/util/dec"
)

var Unmarshaler = dec.Unmarshalers{
	"json": json.Unmarshal,
}
