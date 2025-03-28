package portal

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	ether "github.com/cryptopunkscc/astrald/mod/ether/src"
	tcp "github.com/cryptopunkscc/astrald/mod/tcp/src"
	"github.com/cryptopunkscc/portal/api/astrald"
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/pkg/config"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"slices"
)

type Config struct {
	Dir            string `yaml:",omitempty"`
	Dirs           `yaml:",omitempty"`
	astrald.Config `yaml:",inline"`
	ApphostAddr    string `yaml:"-"`
}

type Dirs struct {
	Astrald  string `yaml:",omitempty"`
	AstralDB string `yaml:",omitempty"`
	Tokens   string `yaml:",omitempty"`
	Apps     string `yaml:",omitempty"`
	Bin      string `yaml:",omitempty"`
}

var baseConfig = Config{
	Dir: "./portal",
	Dirs: Dirs{
		Astrald:  "astrald",
		AstralDB: "astrald",
		Tokens:   "tokens",
		Apps:     "apps",
		Bin:      "bin",
	},
	Config: astrald.Config{
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
	},
}

func (c *Config) dirs() []*string {
	return []*string{
		&c.Astrald,
		&c.AstralDB,
		&c.Dir,
		&c.Tokens,
		&c.Apps,
		&c.Bin,
	}
}

func (c *Config) Yaml() (s string) {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (c *Config) Build() (err error) {
	envConfig := Config{}
	envConfig.readEnvVars()
	config.Merge(c, &envConfig, &baseConfig)
	c.fixPath()
	if err = c.mkdirAll(); err != nil {
		return
	}
	c.fixApphostAddress()
	c.writeEnvVars()
	return
}

func (c *Config) readEnvVars() {
	for e, s := range c.env() {
		*s = e.Get()
	}
}

func (c *Config) writeEnvVars() {
	for e, s := range c.env() {
		e.Set(*s)
	}
}

func (c *Config) env() map[env.Key]*string {
	return map[env.Key]*string{
		env.AstraldHome:   &c.Astrald,
		env.AstraldDb:     &c.AstralDB,
		env.ApphostAddr:   &c.ApphostAddr,
		env.PortaldHome:   &c.Dir,
		env.PortaldTokens: &c.Tokens,
		env.PortaldApps:   &c.Apps,
		env.PortaldBin:    &c.Bin,
	}
}

func (c *Config) fixPath() {
	for _, path := range c.dirs() {
		if filepath.IsLocal(*path) {
			*path = filepath.Join(c.Dir, *path)
		}
	}
}

func (c *Config) mkdirAll() (err error) {
	plog.TraceErr(&err)
	for _, s := range c.dirs() {
		if err = os.MkdirAll(*s, 0755); err != nil {
			return
		}
	}
	return
}

func (c *Config) Load(path ...string) (err error) {
	l := config.Loader[*Config]{
		Unmarshal: yaml.Unmarshal,
		Config:    c,
		File:      DefaultConfigFile,
	}
	if err = l.Load(path...); err != nil {
		return
	}
	if filepath.IsLocal(c.Dir) {
		c.Dir = filepath.Join(l.Dir, c.Dir)
	}
	return
}
func (c *Config) fixApphostAddress() {
	if len(c.ApphostAddr) > 0 {
		c.Apphost.Listen = slices.DeleteFunc(c.Apphost.Listen, func(s string) bool { return s == c.ApphostAddr })
		c.Apphost.Listen = slices.Insert(c.Apphost.Listen, 0, c.ApphostAddr)
	}

	if len(c.Apphost.Listen) > 0 {
		c.ApphostAddr = c.Apphost.Listen[0]
	}
}

const DefaultConfigFile = ".portal.env.yml"
