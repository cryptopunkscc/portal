package astrald

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"reflect"
	"strings"
)

type Config struct {
	Node    core.Config    `yaml:",omitempty"`
	Apphost apphost.Config `yaml:",omitempty"`
	Ether   ether.Config   `yaml:",omitempty"`
	TCP     tcp.Config     `yaml:",omitempty"`
}

func (c Config) Map() (out map[string]any) {
	out = map[string]any{}
	v := reflect.ValueOf(c)
	t := v.Type()
	for ii := range t.NumField() {
		n := t.Field(ii).Name
		n = strings.ToLower(n)
		out[n] = v.Field(ii).Interface()
	}
	return
}
