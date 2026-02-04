package all

import (
	"github.com/cryptopunkscc/portal/pkg/util/dec"
	"github.com/cryptopunkscc/portal/pkg/util/dec/json"
	"github.com/cryptopunkscc/portal/pkg/util/dec/yaml"
)

var Unmarshalers = dec.From(json.Unmarshaler, yaml.Unmarshalers)
