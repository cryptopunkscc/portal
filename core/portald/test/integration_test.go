package portald

import (
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"testing"
	"time"
)

func TestService_Integration(t *testing.T) {
	plog.Verbosity = plog.Debug
	it := newIntegrationTest()
	it.s1.cleanDir(t)
	it.s2.cleanDir(t)
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
			name: "start",
			test: func(t *testing.T) {
				ctx := testServiceContext(t)

				it.s1.configure(t)
				it.s1.testNodeStart(t, ctx)
				it.s1.testNodeAlias(t)
				it.s1.testCreateUser(t)

				t.Run("basic", func(t *testing.T) {
					it.s1.testInstallApps(t)
					it.s1.testUninstallApp(t)
					it.s1.testPublishAppBundle(t)
					it.s1.testReconnectAsUser(t)
					time.Sleep(2000 * time.Millisecond)
					it.s1.awaitPublishedObjects(t)
					it.s1.testSearchObjects(t, "app.manifest")
					it.s1.testFetchReleases(t)
					it.s1.testFetchAppBundleExecs(t)
				})

				t.Run("claim", func(t *testing.T) {
					it.s2.configure(t)
					it.s2.testNodeStart(t, ctx)

					it.s1.testAddEndpoint(t, &it.s2)
					it.s1.testUserClaim(t, &it.s2)
				})

				t.Run("WIP", func(t *testing.T) {
					t.SkipNow() // FIXME
					it.s2.testAddEndpoint(t, &it.s1)
					it.s2.testSearchObjects(t, "app.manifest")
				})
			},
		},
		{
			name: "basic js",
			test: func(t *testing.T) {
				pkg := "test.basic.js"
				ctx := testServiceContext(t)
				s := it.s1
				s.config.Apps = target.Abs("./apps")
				s.configure(t)
				s.testNodeStart(t, ctx)
				s.testSetupToken(t, pkg)
				s.testOpen(t, ctx, pkg)
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
