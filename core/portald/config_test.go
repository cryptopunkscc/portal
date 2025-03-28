package portald

import (
	"github.com/cryptopunkscc/portal/core/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
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
	err := c.build()
	if err != nil {
		plog.Println(err)
		return
	}
	println(c.Yaml())
}
