package portald

import (
	"github.com/cryptopunkscc/astrald/core"
	apphost "github.com/cryptopunkscc/astrald/mod/apphost/src"
	"github.com/cryptopunkscc/portal/core/astrald"
	"github.com/cryptopunkscc/portal/core/env"
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

type AstraldConfigs struct {
	Node    core.Config    `yaml:",omitempty"`
	Apphost apphost.Config `yaml:",omitempty"`
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
	defer plog.TraceErr(&err)
	p, err := getPortalConfigPath(path...)
	if err != nil {
		return
	}
	bytes, err := os.ReadFile(p)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(bytes, c); err != nil {
		return
	}
	return
}

func (c *Config) Save(path ...string) (err error) {
	defer plog.TraceErr(&err)
	p, err := getPortalConfigPath(path...)
	if err != nil {
		return
	}
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	if err = os.WriteFile(p, bytes, 0644); err != nil {
		return
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

func getPortalConfigPath(path ...string) (p string, err error) {
	p = filepath.Join(path...)
	if len(p) == 0 {
		p = DefaultConfigFile
	} else {
		ok := false
		if ok, err = isDir(p); err != nil {
			return
		}
		if ok {
			p = filepath.Join(p, DefaultConfigFile)
		}
	}
	return
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

const DefaultConfigFile = "portal.env.yml"
