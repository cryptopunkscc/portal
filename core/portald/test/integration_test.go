package portald

import (
	"fmt"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
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
			Name: "should start node",
			Test: it.s1.start(),
		},
		{
			Name:    "should run basic js app",
			Test:    it.s1.openApp("test.basic.js"),
			Require: test.Tests{it.s1.setupToken("test.basic.js")},
		},
		{
			Name: "should test core js lib",
			Require: test.Tests{
				buildCoreJsTestCommon(),
				it.s1.start(),
				it.s1.installApps("/core/js/test/common"),
				it.s1.openApp("portal.js.test.common"),
			},
		},
		{
			Name:    "should get alias",
			Test:    it.s1.nodeAlias(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name: "should create user",
			Test: it.s1.createUser(),
		},
		{
			Name: "should add endpoint",
			Test: it.s1.addEndpoint(&it.s2),
		},
		{
			Name: "should claim",
			Test: it.s1.userClaim(&it.s2),
		},
		//{
		//	Name:    "should build apps",
		//	Test:    it.s1.buildApps(),
		//	Require: test.Tests{it.s1.start()},
		//},
		{
			Name:    "should install apps",
			Test:    it.s1.installDefaultApps(),
			Require: test.Tests{it.s1.start()},
		},
		{
			Name:    "should uninstall app",
			Test:    it.s1.uninstallApp(),
			Require: test.Tests{it.s1.installDefaultApps()},
		},
		{
			Name:    "should publish app bundles",
			Test:    it.s1.publishAppBundles(),
			Require: test.Tests{it.s1.createUser()},
		},
		{
			Name:    "should await published app bundles",
			Test:    it.s1.awaitPublishedBundles(),
			Require: test.Tests{it.s1.publishAppBundles(), it.s1.createUser()},
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
			Name:    "should scan its own objects",
			Test:    it.s1.scanObjects("app.manifest"),
			Require: test.Tests{it.s1.awaitPublishedBundles()},
		},
		{
			Name: "should scan another node's objects",
			Test: it.s2.scanObjects("app.manifest", &it.s1),
			Require: test.Tests{
				it.s1.scanObjects("app.manifest"),
				it.s1.userClaim(&it.s2),
				it.s2.addEndpoint(&it.s1),
			},
		},
		{
			Name:    "should has user",
			Test:    it.s1.hasUser(),
			Require: test.Tests{it.s1.createUser()},
		},
		{
			Name:    "should has user",
			Test:    it.s2.hasUser(),
			Require: test.Tests{it.s1.userClaim(&it.s2)},
		},
		{
			Name:    "should list siblings",
			Test:    it.s1.listSiblings(),
			Require: test.Tests{it.s1.userClaim(&it.s2)},
		},
		{
			Name:    "should list siblings 2",
			Test:    it.s2.listSiblings(),
			Require: test.Tests{it.s1.listSiblings()},
		},
		{
			Name: "should list available apps 1",
			Test: it.s1.availableApps(),
			Require: test.Tests{
				it.s1.awaitPublishedBundles(),
			},
		},
		{
			Name: "should list available apps 2",
			Test: it.s2.availableApps(),
			Require: test.Tests{
				it.s1.awaitPublishedBundles(),
				it.s1.userClaim(&it.s2),
				it.s2.addEndpoint(&it.s1),
			},
		},
		{
			Name:    "should install apps by package name",
			Test:    it.s2.installAppsByPackage(&it.s1),
			Require: test.Tests{it.s2.availableApps()},
		},
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
	i.s1.name = "s1"
	i.s1.config.ApphostAddr = "tcp:127.0.0.1:8636"
	i.s1.config.Apphost.Listen = []string{"tcp:127.0.0.1:8636"}
	i.s1.config.Apphost.ObjectServer.Bind = []string{"tcp:127.0.0.1:8625"}
	i.s1.config.Ether.UDPPort = 8833
	i.s1.config.TCP.ListenPort = 1796

	i.s2.ctx = ctx
	i.s2.name = "s2"
	i.s2.config.ApphostAddr = "tcp:127.0.0.1:8637"
	i.s2.config.Apphost.Listen = []string{"tcp:127.0.0.1:8637"}
	i.s2.config.Apphost.ObjectServer.Bind = []string{"tcp:127.0.0.1:8626"}
	i.s2.config.Ether.UDPPort = 8834
	i.s2.config.TCP.ListenPort = 1797

	return
}
