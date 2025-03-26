package portald

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/config"
	"os"
	"path/filepath"
)

type Config struct {
	Dir string
	Dirs
	AstraldConfigs `yaml:",inline"`
	ApphostAddr    string `yaml:"-"`
}

type Dirs struct {
	Astrald  string
	AstralDB string
	Tokens   string
	Apps     string
	Bin      string
}

type AstraldConfigs struct {
	Node    core.Config
	Apphost apphost.Config
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
	AstraldConfigs: AstraldConfigs{
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

func (c *Config) build() *Config {
	envConfig := Config{}
	envConfig.readEnvVars()
	config.Merge(c, &envConfig, &baseConfig)
	c.fixPath()
	c.mkdirAll()
	c.writeEnvVars()
	return c
}

func (c *Config) readEnvVars() {
	for e, s := range c.env() {
		*s = e.Get()
	}
	if len(c.ApphostAddr) > 0 {
		c.Apphost.Listen = append([]string{c.ApphostAddr}, c.Apphost.Listen...)
	}
}

func (c *Config) writeEnvVars() {
	if len(c.Apphost.Listen) > 0 {
		c.ApphostAddr = c.Apphost.Listen[0]
	}
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

func (c *Config) mkdirAll() {
	for _, s := range c.dirs() {
		if err := os.MkdirAll(*s, 0755); err != nil {
			panic(err)
		}
	}
}
