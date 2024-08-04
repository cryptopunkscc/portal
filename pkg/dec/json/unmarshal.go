package json

import (
	"encoding/json"
	"github.com/cryptopunkscc/portal/pkg/dec"
)

var Unmarshaler = dec.Unmarshalers{
	"json": json.Unmarshal,
}
