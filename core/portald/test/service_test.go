package portald

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/nodes"
	"github.com/cryptopunkscc/portal/api/portal"
	portald2 "github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/astrald/debug"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/portald"
	golang "github.com/cryptopunkscc/portal/pkg/go"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/runner/goja"
	source2 "github.com/cryptopunkscc/portal/source"
	"github.com/cryptopunkscc/portal/source/app"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/npm"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
)

type testService struct {
	name   string
	alias  string
	ctx    context.Context
	config portal.Config
	*portald.Service
	apps       []target.Portal_
	published  map[astral.ObjectID]bundle.Info
	published2 map[astral.ObjectID]app.ReleaseInfo
}

func testServiceContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
	})
	return ctx
}

func (s *testService) cleanDir(t *testing.T) {
	s.config.Dir = test.CleanDir(t, ".test_"+s.name)
}

func (s *testService) test(run func(t *testing.T)) test.Test {
	name := fmt.Sprintf("%s.%s", s.name, callerName())
	return test.New(name, run)
}

func (s *testService) arg(name string) testBuilder {
	return testBuilder{
		testService: s,
		Suffix:      name,
	}
}

type testBuilder struct {
	*testService
	Suffix string
}

func (s testBuilder) test(run func(t *testing.T)) test.Test {
	name := fmt.Sprintf("%s.%s(%s)", s.name, callerName(), s.Suffix)
	return test.New(name, run)
}

func callerName() (name string) {
	if pc, _, _, ok := runtime.Caller(2); ok {
		if funcObj := runtime.FuncForPC(pc); funcObj != nil {
			name = funcObj.Name()
			c := strings.Split(name, ".")
			name = c[len(c)-1]
		}
	}
	return
}

//go:embed apps
var AppsFS embed.FS

func (s *testService) configure() test.Test {
	return s.test(func(t *testing.T) {
		s.Service = &portald.Service{}
		s.Config = s.config
		s.Config.Node.Log.Level = 100
		s.ExtraTokens = []string{"portal"}
		s.AppSources = []target.Source{source.Embed(AppsFS)}
		err := s.Configure()
		test.AssertErr(t, err)
		//s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald} // Faster testing
		s.Astrald = &debug.Astrald{NodeRoot: s.Config.Astrald} // Debugging astrald

		f := bind.AutoTokenCoreFactory{Tokens: s.Tokens()}
		s.Resolve = target.Any[target.Runnable](
			target.Skip("node_modules"),
			goja.Runner(f.Create2).Try,
		)
		//s.Apphost.Log = plog.New().Scope("TEST_APPHOST_" + s.name)
	})
}

func (s *testService) start() test.Test {
	return s.test(func(t *testing.T) {
		err := s.Start(s.ctx)
		test.AssertErr(t, err)
	}).Requires(
		s.configure(),
	)
}

func (s *testService) nodeAlias() test.Test {
	return s.test(func(t *testing.T) {
		alias, err := s.Apphost.NodeAlias()
		test.AssertErr(t, err)
		assert.NotZero(t, alias)
		s.alias = alias
	})
}

func (s *testService) createUser() test.Test {
	return s.test(func(t *testing.T) {
		err := s.CreateUser("test_user")
		test.AssertErr(t, err)
	}).Requires(
		s.start(),
	)
}

func (s *testService) userInfo() test.Test {
	return s.test(func(t *testing.T) {
		info, err := s.UserInfo()
		test.AssertErr(t, err)
		assert.NotZero(t, info)
		plog.Println(*info)
	})
}

func (s *testService) hasUser() test.Test {
	return s.test(func(t *testing.T) {
		b := s.HasUser()
		assert.True(t, b)
	})
}

func (s *testService) userClaim(s2 *testService) test.Test {
	return s.test(func(t *testing.T) {
		err := s.Claim(s2.Apphost.TargetID.String())
		test.AssertErr(t, err)
	}).Requires(
		s.createUser(),
		s.addEndpoint(s2),
	)
}

func (s *testService) addEndpoint(s2 *testService) test.Test {
	return s.arg(s2.name).test(func(t *testing.T) {
		id := s2.Apphost.TargetID.String()
		port := s2.Config.TCP.ListenPort
		endpoint := fmt.Sprintf("tcp:127.0.0.1:%d", port)
		err := nodes.Op(&s.Apphost).AddEndpoint(id, endpoint)
		test.AssertErr(t, err)
	}).Requires(
		s.start(),
		s2.start(),
	)
}

func (s *testService) buildApps() test.Test {
	return s.test(func(t *testing.T) {
		err := apps.Build("pack")
		if err.Error() != "npm is required but not installed" {
			test.AssertErr(t, err)
		}
	})
}

func (s *testService) installApps(path string) test.Test {
	return s.test(func(t *testing.T) {
		r, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		l, err := source.File(filepath.Join(r, path))
		assert.NoError(t, err)
		ctx := context.Background()
		for _, r := range s.Installer().Dispatcher().List(l) {
			err := r.Run(ctx)
			test.AssertErr(t, err)
			s.apps = append(s.apps, r)
			break
		}
		assert.NotEmpty(t, s.apps)
	})
}

func (s *testService) installDefaultApps() test.Test {
	return s.test(func(t *testing.T) {
		l := source.Embed(apps.Builds)
		ctx := context.Background()
		for _, r := range s.Installer().Dispatcher().List(l) {
			err := r.Run(ctx)
			test.AssertErr(t, err)
			s.apps = append(s.apps, r)
			break
		}
		assert.NotEmpty(t, s.apps)
	})
}

func (s *testService) installAppsByPackage(s2 *testService) test.Test {
	return s.test(func(t *testing.T) {
		ctx := context.Background()
		d := s.Installer().Dispatcher()
		count := 0
		for _, info := range s2.published {
			err := d.Run(ctx, info.Manifest.Package)
			assert.NoError(t, err)
			count++
			break
		}
		assert.NotEmpty(t, count)
	})
}

func (s *testService) uninstallApp() test.Test {
	return s.test(func(t *testing.T) {
		if len(s.apps) == 0 {
			test.AssertErr(t, plog.Errorf("no apps installed"))
		}
		p := s.apps[0].Manifest().Package
		err := s.Installer().Uninstall(p)
		test.AssertErr(t, err)
	})
}

func (s *testService) publishAppBundles() test.Test {
	return s.test(func(t *testing.T) {
		b := source.Embed(apps.Builds)
		s.published = map[astral.ObjectID]bundle.Info{}

		l, err := s.PublishAppsFS(b)
		test.AssertErr(t, err)
		count := 0
		for _, app := range l {
			s.published[*app.ReleaseID] = app
			count++
		}
		assert.NotZero(t, count)
	})
}

func (s *testService) awaitPublishedBundles() test.Test {
	return s.test(func(t *testing.T) {
		if len(s.published) == 0 {
			t.Fatalf("no published bundles")
		}
		for id := range s.published {
			s.awaitObject(id).Run(t)
			break
		}
	})
}

func (s *testService) publishAppBundlesV2() test.Test {
	return s.test(func(t *testing.T) {
		s.published2 = map[astral.ObjectID]app.ReleaseInfo{}
		src := source2.FSRef(apps.Builds)
		l, err := app.PublishAppBundlesSrc(s.Apphost.Client, src)
		test.AssertErr(t, err)
		count := 0
		for _, app := range l {
			s.published2[*app.ReleaseID] = app
			count++
		}
		assert.NotZero(t, count)
	})
}

func (s *testService) awaitPublishedBundlesV2() test.Test {
	return s.test(func(t *testing.T) {
		if len(s.published2) == 0 {
			t.Fatalf("no published bundles")
		}
		for id := range s.published2 {
			s.awaitObject(id).Run(t)
			break
		}
	})
}

func (s *testService) awaitObject(id astral.ObjectID) test.Test {
	return s.test(func(t *testing.T) {
		c := s.Apphost.Objects()
		limit := 10
		for {
			rc, err := c.Read(nil, &id, 0, 0)
			if err != nil {
				if limit > 0 {
					limit--
					time.Sleep(200 * time.Millisecond)
					continue
				}
				test.AssertErr(t, err)
			}

			buf := bytes.NewBuffer(nil)
			_, err = buf.ReadFrom(rc)
			test.AssertErr(t, err)

			plog.Println(id, buf.String())
			break
		}
	})
}

func (s *testService) readObject(id astral.ObjectID) test.Test {
	return s.test(func(t *testing.T) {
		c := s.Apphost.Objects()
		rc, err := c.Read(nil, &id, 0, 0)
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		test.AssertErr(t, err)

		plog.Println(id, buf.String())
	})
}

func (s *testService) reconnectAsUser() test.Test {
	return s.test(func(t *testing.T) {
		s.Apphost.Token = s.UserCreated.AccessToken.String()
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) reconnectAsUser2(s2 *testService) test.Test {
	return s.test(func(t *testing.T) {
		s.Apphost.Token = s2.UserCreated.AccessToken.String()
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) reconnectAs(alias string) test.Test {
	return s.arg(alias).test(func(t *testing.T) {
		pt, err := s.Tokens().Resolve(alias)
		test.AssertErr(t, err)
		s.Apphost.Token = pt.Token.String()
		err = s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) scanObjects(s2 ...*testService) test.Test {
	o := s2
	if len(o) == 0 {
		o = append(o, s)
	}
	return s.test(func(t *testing.T) {
		scan, err := s.Apphost.Objects().Scan(nil, "", false)
		test.AssertErr(t, err)

		count := 0
		for result := range scan {
			count++
			plog.Println(result)
		}
		assert.Greater(t, count, 0)
	})
}

func (s *testService) availableApps() test.Test {
	return s.test(func(t *testing.T) {
		aa, err := s.AvailableApps(s.ctx, false)
		test.AssertErr(t, err)

		count := 0
		for app := range aa {
			count++
			plog.Println("available app:", app)
		}
		assert.NotZero(t, count)
	})
}

func (s *testService) searchObjects(query string, s2 ...*testService) test.Test {
	o := s2
	if len(o) == 0 {
		o = append(o, s)
	}
	return s.arg(query).test(func(t *testing.T) {
		results, err := s.Apphost.Objects().Search(nil, query)
		test.AssertErr(t, err)
		count := 0
		for result := range results {
			count++
			plog.Println(result)
		}
		assert.NotZero(t, count)
	})
}

func (s *testService) fetchReleases() test.Test {
	return s.test(func(t *testing.T) {
		for id, info := range s.published {
			r := &bundle.Release{}
			err := s.Apphost.Objects().Fetch(&id, r)
			test.AssertErr(t, err)
			assert.Equal(t, *info.Release.BundleID, *r.BundleID)
			assert.Equal(t, *info.Release.ManifestID, *r.ManifestID)
			assert.Equal(t, info.Release.Release, r.Release)
			assert.Equal(t, info.Release.Target, r.Target)
		}
	})
}

func (s *testService) fetchReleasesV2() test.Test {
	return s.test(func(t *testing.T) {
		for id, info := range s.published2 {
			obj, err := s.Apphost.Objects().Get(&id)
			r := obj.(*app.ReleaseMetadata)
			//a, err := app.Objects{Adapter: &s.Apphost}.GetAppBundle(info.Manifest.Package)
			test.NoError(t, err)
			assert.Equal(t, *info.BundleID, *r.BundleID)
			assert.Equal(t, *info.ManifestID, *r.ManifestID)
			assert.Equal(t, info.Release, r.Release)
			assert.Equal(t, info.Target, r.Target)
		}
	})
}

func (s *testService) signAppContract(pkg string) test.Test {
	return s.arg(pkg).test(func(t *testing.T) {
		id, err := s.Apphost.Resolve(pkg)
		test.AssertErr(t, err)
		contract, err := s.Apphost.SignAppContract(id)
		test.AssertErr(t, err)
		plog.Println(contract)
	})
}

func (s *testService) fetchAppBundleExecs() test.Test {
	return s.test(func(t *testing.T) {
		for _, r := range s.published {
			_, err := s.Bundles().GetByObjectID(*r.Release.BundleID)
			test.AssertErr(t, err)
		}
	})
}

func (s *testService) openApp(pkg string) test.Test {
	return s.arg(pkg).test(func(t *testing.T) {
		o := portald2.OpenOpt{}
		err := s.Open().Run(s.ctx, o, pkg)
		test.AssertErr(t, err)
	})
}

func (s *testService) setupToken(pkg string) test.Test {
	return s.arg(pkg).test(func(t *testing.T) {
		_, err := s.Tokens().Resolve(pkg)
		test.AssertErr(t, err)
	}).Requires(
		s.start(),
	)
}

func (s *testService) listSiblings() test.Test {
	return s.test(func(t *testing.T) {
		siblings, err := s.Apphost.User().Siblings(nil)
		test.AssertErr(t, err)
		count := 0
		for sibling := range siblings {
			count++
			plog.Println(sibling.String())
		}
		assert.Equal(t, 1, count)
	})
}

func buildCoreJsTestCommon() test.Test {
	return test.New(test.CallerName(), func(t *testing.T) {
		root, err := golang.FindProjectRoot()
		assert.NoError(t, err)
		d := source.Dir(root, "core", "js", "test", "common")
		err = npm.BuildRunner().Run(context.TODO(), d)
		assert.NoError(t, err)
	})
}
