package portald

import (
	"fmt"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"testing"
)

func TestService_Integrations(t *testing.T) {
	plog.Verbosity = plog.Debug
	it := newIntegrationTest(t)
	it.s1.cleanDir(t)
	it.s2.cleanDir(t)
	runner := test.Runner{}
	tests := []test.Task{
		{
			Name: "should configure portald",
			Test: it.s1.configure(),
		},
		{
			Name:    "should start node",
			Test:    it.s1.start(),
			Require: test.Tests{it.s1.configure()},
		},
		{
			Name:    "should get alias",
			Test:    it.s1.nodeAlias(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name:    "should create user",
			Test:    it.s1.createUser(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name:    "should add endpoint",
			Test:    it.s1.addEndpoint(&it.s2),
			Require: test.Tests{it.s1.start(), it.s2.start()},
		},
		{
			Name:    "should claim",
			Test:    it.s1.userClaim(&it.s2),
			Require: test.Tests{it.s1.createUser(), it.s1.addEndpoint(&it.s2)},
		},
		//{
		//	Name:    "should build apps",
		//	Test:    it.s1.buildApps(),
		//	Require: test.Tests{it.s1.start()},
		//},
		{
			Name:    "should install apps",
			Test:    it.s1.installApps(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name:    "should uninstall app",
			Test:    it.s1.uninstallApp(),
			Require: test.Tests{it.s1.installApps()},
		},
		{
			Name:    "should publish app bundles",
			Test:    it.s1.publishAppBundles(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name:    "should reconnect as user",
			Test:    it.s1.reconnectAsUser(),
			Require: test.Tests{it.s1.createUser()},
		},
		{
			Name:    "should await published app bundles",
			Test:    it.s1.awaitPublishedBundles(),
			Require: test.Tests{it.s1.publishAppBundles(), it.s1.reconnectAsUser()},
		},
		{
			Name:    "should search objects",
			Test:    it.s1.searchObjects("app.manifest"),
			Require: test.Tests{it.s1.awaitPublishedBundles()},
		},
		{
			Name:    "should fetch releases",
			Test:    it.s1.fetchReleases(),
			Require: test.Tests{it.s1.awaitPublishedBundles()},
		},
		{
			Name:    "should fetch executable app bundles",
			Test:    it.s1.fetchAppBundleExecs(),
			Require: test.Tests{it.s1.awaitPublishedBundles()},
		},
		{
			Name:    "should run basic js app",
			Test:    it.s1.openApp("test.basic.js"),
			Require: test.Tests{it.s1.setupToken("test.basic.js")},
		},
		{
			Name:    "should search object from parent node WIP",
			Test:    it.s2.searchObjects("app.manifest"),
			Require: test.Tests{it.s1.awaitPublishedBundles(), it.s2.addEndpoint(&it.s1)},
		},

		{Test: it.s1.setupToken("test.basic.js"), Require: test.Tests{it.s1.start()}},
		{Test: it.s2.start(), Require: test.Tests{it.s2.configure()}},
		{Test: it.s2.addEndpoint(&it.s1), Require: test.Tests{it.s1.start(), it.s2.start()}},
		{Test: it.s2.configure()},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d  %s", i, tt.Name), runner.Run(tests, tt))
	}
}

type integrationTest struct {
	s1 testService
	s2 testService
}

func newIntegrationTest(t *testing.T) (i integrationTest) {
	ctx := testServiceContext(t)

	i.s1.ctx = ctx
	i.s1.name = ".test1"
	i.s1.config.ApphostAddr = "tcp:127.0.0.1:8636"
	i.s1.config.Apphost.Listen = []string{"tcp:127.0.0.1:8636"}
	i.s1.config.Ether.UDPPort = 8833
	i.s1.config.TCP.ListenPort = 1796

	i.s2.ctx = ctx
	i.s2.name = ".test2"
	i.s2.config.ApphostAddr = "tcp:127.0.0.1:8637"
	i.s2.config.Apphost.Listen = []string{"tcp:127.0.0.1:8637"}
	i.s2.config.Ether.UDPPort = 8834
	i.s2.config.TCP.ListenPort = 1797

	return
}
