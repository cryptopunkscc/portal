package main

import (
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
)

func TestApplication_loadConfig_custom(t *testing.T) {
	unsetEnv()
	dir := test.Dir(t)
	test.Clean(dir)
	dir = test.Mkdir(t)
	c := portal.Config{}
	c.Dir = dir
	c.Node.Log.Level = 100
	c.Apphost.Listen = []string{"tcp:127.0.0.1:8635"}
	if err := writeConfig(c, dir, portal.DefaultConfigFile); err != nil {
		plog.P().Println(err)
	}

	args := RunArgs{ConfigPath: dir}
	err := application.loadConfig(args)
	if err != nil {
		plog.P().Println(err)
	}
	assert.Equal(t, c, application.Config)
}

func TestApplication_loadConfig_platformDefault(t *testing.T) {
	if err := application.loadConfig(RunArgs{}); err != nil {
		plog.P().Println(err)
	}
	if err := application.Configure(); err != nil {
		plog.P().Println(err)
	}
}

func unsetEnv() {
	for _, key := range []env.Key{
		env.AstraldHome,
		env.AstraldDb,
		env.ApphostAddr,
		env.PortaldHome,
		env.PortaldTokens,
		env.PortaldApps,
		env.PortaldBin,
	} {
		key.Unset()
	}
}

func writeConfig(c portal.Config, path ...string) (err error) {
	defer plog.TraceErr(&err)
	p := filepath.Join(path...)
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	if err = os.WriteFile(p, bytes, 0644); err != nil {
		return
	}
	return
}
