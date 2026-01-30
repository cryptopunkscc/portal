package debug

import (
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/user"
	"github.com/cryptopunkscc/portal/cmd/portal-goja/src"
	"github.com/cryptopunkscc/portal/core/apphost"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/runner/astrald/debug"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/source/npm"
	"github.com/cryptopunkscc/portal/source/tmpl"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {

	ctx := TestContext{
		Context: astral.NewContext(nil),
		Name:    "1",
		Astrald: debug.Astrald{},
	}
	ctx.Astrald.NodeRoot = test.CleanMkdir(t, ".test_"+ctx.Name)

	js := &TestApp{
		Template: "js",
		Name:     "js-app/js-app",
		Type:     "js",
	}
	jsRollup := &TestApp{
		Template: "js-rollup",
		Name:     "js-rollup-app",
		Type:     "js",
	}
	svelte := &TestApp{
		Template: "svelte",
		Name:     "svelte-app",
		Type:     "html",
	}

	projectRoot, err := golang.FindProjectRoot()
	test.NoError(t, err)
	coreJsTest := &TestApp{
		Name: "core.js.test",
		Type: "js",
		Path: path.Join(projectRoot, "core/js/test/common"),
	}

	runner := test.Runner{}
	tests := []test.Task{
		{
			Name: "start astrald",
			Test: ctx.StartAstrald(),
		},
		{
			Name: "create user",
			Test: ctx.CreateUser(),
		},
		{
			Name: "user info",
			Test: ctx.GetUserInfo(),
		},
		{
			Name: "create js project",
			Test: ctx.CreateProject(jsRollup),
		},
		{
			Name: "build js project",
			Test: ctx.BuildProject(jsRollup),
		},
		{
			Name: "run js app by path",
			Test: ctx.RunAppByPath(jsRollup),
		},
		{
			Name: "publish js app",
			Test: ctx.PublishApp(jsRollup),
		},
		{
			Name: "run js app by name",
			Test: ctx.RunAppByName(jsRollup),
		},
		{
			Name: "run js app by package",
			Test: ctx.RunAppByPackage(jsRollup),
		},
		{
			Name: "run js app by release ID",
			Test: ctx.RunAppByReleaseID(jsRollup),
		},
		{
			Name: "run js app by bundle ID",
			Test: ctx.RunAppByBundleID(jsRollup),
		},
		{
			Name: "publish js app",
			Test: ctx.RunAppByReleaseID(js),
		},
		{
			Name: "publish svelte app",
			Test: ctx.PublishApp(svelte),
		},
		{
			Name:    "test core js",
			Test:    ctx.RunAppByPath(coreJsTest),
			Require: test.Tests{ctx.CreateUser()},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d  %s", i, tt.Name), runner.Run(tests, tt))
	}
}

type TestContext struct {
	*astral.Context
	Name           string
	Astrald        debug.Astrald
	Apphost        apphost.Adapter
	UserCreateInfo *user.CreatedUserInfo
	UserInfo       *user.Info
}

type TestApp struct {
	Template    string
	Name        string
	Type        string
	Path        string
	ReleaseInfo *app.ReleaseInfo
}

func (a TestApp) GetPath(ctx *TestContext) (projectPath string) {
	if a.Path == "" {
		a.Path = path.Join(ctx.Astrald.NodeRoot, a.Name)
	}
	if projectPath = a.Path; projectPath == "" {
		projectPath = path.Join(ctx.Astrald.NodeRoot, a.Name)
	}
	return
}

func (c *TestContext) Test() test.Test {
	name := fmt.Sprintf("%s-%s", c.Name, test.CallerName(2))
	return test.New(name, func(t *testing.T) {})
}

func (c *TestContext) StartAstrald() test.Test {
	return c.Test().Func(func(t *testing.T) {
		err := c.Astrald.Start(c)
		test.NoError(t, err)

		time.Sleep(time.Second * 1)

		err = c.Apphost.Connect()
		test.NoError(t, err)
	})
}

func (c *TestContext) CreateUser() test.Test {
	return c.Test().Func(func(t *testing.T) {
		var err error
		c.UserCreateInfo, err = c.Apphost.User().Create(c.Context, "test_user")
		test.NoError(t, err)
		t.Log(c.UserCreateInfo)

		c.Apphost.Token = c.UserCreateInfo.AccessToken.String()
		err = c.Apphost.Reconnect()
		test.NoError(t, err)
	}).Requires(
		c.StartAstrald(),
	)
}

func (c *TestContext) GetUserInfo() test.Test {
	return c.Test().Func(func(t *testing.T) {
		var err error
		c.UserInfo, err = c.Apphost.User().Info(c.Context)
		test.NoError(t, err)
		t.Log(c.UserInfo)
	}).Requires(
		c.CreateUser(),
	)
}

func (c *TestContext) CreateProject(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		p := testApp.GetPath(c)
		if tmpl.CheckTargetDirNotExist(p) != nil {
			return
		}
		err := tmpl.Create(testApp.Template, p)
		test.NoError(t, err)
	})
}

func (c *TestContext) BuildProject(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		err := npm.BuildNpmApps(
			npm.BuildNpmAppsOpt{Pack: true},
			testApp.GetPath(c),
		)
		test.NoError(t, err)
	}).Requires(
		c.CreateProject(testApp),
	)
}

func (c *TestContext) RunAppByPath(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		c.runApp(t, testApp, path.Join(testApp.GetPath(c), "dist"))
	}).Requires(
		c.BuildProject(testApp),
	)
}

func (c *TestContext) PublishApp(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		publisher := app.Publisher{ObjectsClient: c.Apphost.Objects()}
		dir := testApp.GetPath(c)
		if testApp.Template == "js" {
			dir = path.Dir(dir)
		}
		bundles, err := publisher.PublishBundles(c.Context, dir)
		test.NoError(t, err)
		require.Len(t, bundles, 1)
		testApp.ReleaseInfo = &bundles[0]
	}).Requires(
		c.CreateUser(),
		c.BuildProject(testApp),
	)
}

func (c *TestContext) RunAppByName(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		c.runApp(t, testApp, testApp.ReleaseInfo.Manifest.Name)
	}).Requires(
		c.PublishApp(testApp),
	)
}

func (c *TestContext) RunAppByPackage(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		c.runApp(t, testApp, testApp.ReleaseInfo.Manifest.Package)
	}).Requires(
		c.PublishApp(testApp),
	)
}

func (c *TestContext) RunAppByBundleID(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		c.runApp(t, testApp, testApp.ReleaseInfo.BundleID.String())
	}).Requires(
		c.PublishApp(testApp),
	)
}

func (c *TestContext) RunAppByReleaseID(testApp *TestApp) test.Test {
	return c.Test().Args(testApp.Name).Func(func(t *testing.T) {
		c.runApp(t, testApp, testApp.ReleaseInfo.ReleaseID.String())
	}).Requires(
		c.PublishApp(testApp),
	)
}

func (c *TestContext) runApp(t *testing.T, testApp *TestApp, src string) {
	var err error = plog.Errorf("unsupported type")
	switch testApp.Type {
	case "html":
		t.Skip("not supported in debuggable test")
	case "js":
		a := portal_goja.Application{Adapter: &c.Apphost}
		err = a.Run(c.Context, src)
	}
	test.NoError(t, err)
}
