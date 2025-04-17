package portald

import (
	"github.com/cryptopunkscc/portal/pkg/plog"
	"testing"
)

func TestService_Integration(t *testing.T) {
	plog.Verbosity = plog.Debug
	it := newIntegrationTest()
	it.s1.setupDir(t)
	it.s2.setupDir(t)
	tests := []struct {
		name string
	}{
		{name: "1"},
		{name: "2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testServiceContext(t)
			it.s1.configure()
			it.s1.testNodeStart(t, ctx)
			it.s1.testNodeAlias(t)

			it.s2.configure()
			it.s2.testNodeStart(t, ctx)
			it.s2.testNodeAlias(t)
		})
	}
}

func newIntegrationTest() (i integrationTest) {
	i.s1.name = ".test1"
	i.s1.config.ApphostAddr = "tcp:127.0.0.1:8636"
	i.s1.config.Apphost.Listen = []string{"tcp:127.0.0.1:8636"}
	i.s1.config.Ether.UDPPort = 8833

	i.s2.name = ".test2"
	i.s2.config.ApphostAddr = "tcp:127.0.0.1:8637"
	i.s2.config.Apphost.Listen = []string{"tcp:127.0.0.1:8637"}
	i.s2.config.Ether.UDPPort = 8834
	return
}

type integrationTest struct {
	s1 testService
	s2 testService
}
