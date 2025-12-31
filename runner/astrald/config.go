package astrald

import (
	"reflect"
	"strings"

	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
)

type Config struct {
	Node    core.Config    `yaml:",omitempty"`
	Apphost apphost.Config `yaml:",omitempty"`
	Ether   ether.Config   `yaml:",omitempty"`
	TCP     tcp.Config     `yaml:",omitempty"`
}

var DefaultConfig = Config{
	Node: core.Config{
		Log: core.LogConfig{
			Level:         100,
			DisableColors: false,
		},
	},
	Apphost: apphost.Config{
		Workers: 32,
		Listen: []string{
			"tcp:127.0.0.1:8625",
			"unix:~/.apphost.sock",
			"memu:apphostu",
			"memb:apphostb",
		},
		ObjectServer: apphost.ObjectServerConfig{
			Bind: []string{
				"tcp:127.0.0.1:8624",
			},
		},
	},
	Ether: ether.Config{
		UDPPort: 8822,
	},
	TCP: tcp.Config{
		ListenPort: 1791,
	},
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
