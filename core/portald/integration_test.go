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
	for _, tt := range []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "just start",
			test: func(t *testing.T) {
				ctx := testServiceContext(t)
				it.s1.configure(t)
				it.s1.testNodeStart(t, ctx)
				it.s1.testNodeAlias(t)
			},
		},
		{
			name: "create user and claim",
			test: func(t *testing.T) {
				ctx := testServiceContext(t)

				it.s2.configure(t)
				it.s2.testNodeStart(t, ctx)
				it.s2.testNodeAlias(t)

				it.s1.configure(t)
				it.s1.testNodeStart(t, ctx)
				it.s1.testNodeAlias(t)
				it.s1.testCreateUser(t)

				it.s1.testAddEndpoint(t, &it.s2)
				it.s1.testUserClaim(t, &it.s2)
			},
		},
	} {
		t.Run(tt.name, tt.test)
	}
}

func newIntegrationTest() (i integrationTest) {
	i.s1.name = ".test1"
	i.s1.config.ApphostAddr = "tcp:127.0.0.1:8636"
	i.s1.config.Apphost.Listen = []string{"tcp:127.0.0.1:8636"}
	i.s1.config.Ether.UDPPort = 8833
	i.s1.config.TCP.ListenPort = 1796

	i.s2.name = ".test2"
	i.s2.config.ApphostAddr = "tcp:127.0.0.1:8637"
	i.s2.config.Apphost.Listen = []string{"tcp:127.0.0.1:8637"}
	i.s2.config.Ether.UDPPort = 8834
	i.s2.config.TCP.ListenPort = 1797
	return
}

type integrationTest struct {
	s1 testService
	s2 testService
}
