package all

import (
	"github.com/cryptopunkscc/portal/pkg/dec"
	"github.com/cryptopunkscc/portal/pkg/dec/json"
	"github.com/cryptopunkscc/portal/pkg/dec/yaml"
)

var Unmarshalers dec.Unmarshaler = dec.From(json.Unmarshaler, yaml.Unmarshalers)
