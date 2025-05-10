package main

import (
	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
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
	c.Apphost.Listen = []string{"tcp:127.0.0.1:8638"}
	err := writeConfig(t, c, dir, portal.DefaultConfigFile)
	test.AssertErr(t, err)

	a := testApplication()
	args := RunArgs{ConfigPath: dir}
	err = a.loadConfig(args)
	test.AssertErr(t, err)
	assert.Equal(t, c, a.Config)
}

func TestApplication_loadConfig_platformDefault(t *testing.T) {
	a := testApplication()
	err := a.loadConfig(RunArgs{})
	test.AssertErr(t, err)

	err = a.Configure()
	test.AssertErr(t, err)
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

func writeConfig(t *testing.T, c portal.Config, path ...string) (err error) {
	defer plog.TraceErr(&err)
	p := filepath.Join(path...)
	bytes, err := yaml.Marshal(c)
	test.AssertErr(t, err)
	err = os.WriteFile(p, bytes, 0644)
	test.AssertErr(t, err)
	return
}
