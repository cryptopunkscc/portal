package portal

import (
	"testing"

	"github.com/cryptopunkscc/portal/api/env"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
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
	println(c.Yaml())
}
