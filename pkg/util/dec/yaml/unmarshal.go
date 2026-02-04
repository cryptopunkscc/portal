package yaml

import (
	"github.com/cryptopunkscc/portal/pkg/util/dec"
	"gopkg.in/yaml.v3"
)

var Unmarshalers = dec.Unmarshalers{
	"yaml": yaml.Unmarshal,
	"yml":  yaml.Unmarshal,
}
