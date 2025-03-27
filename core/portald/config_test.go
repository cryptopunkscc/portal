package portald

import (
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"gopkg.in/yaml.v3"
	"testing"
)

func init() {
	plog.Verbosity = 100
}

func TestConfig_Build(t *testing.T) {
	env.PortaldBin.Set("envbin")
	c := Config{}
	c.Dir = test.Dir(t, ".test", "portal")
	c.Apps = "argapps"
	c.ApphostAddr = "localhost"
	err := c.Build()
	if err != nil {
		plog.Println(err)
		return
	}
	bytes, err := yaml.Marshal(c)
	if err != nil {
		t.Error(err)
	}
	println(string(bytes))
}
