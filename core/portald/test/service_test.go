package portald

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/object"
	"github.com/cryptopunkscc/portal/api/nodes"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/portal"
	portald2 "github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/astrald/debug"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testService struct {
	name   string
	alias  string
	ctx    context.Context
	config portal.Config
	*portald.Service[target.Portal_]
	apps      []target.Portal_
	published map[object.ID]bundle.Release
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
	s.config.Dir = test.CleanDir(t, s.name)
}

func (s *testService) test(name string, run func(t *testing.T)) test.Test {
	return test.New(s.name+" "+name, run)
}

//go:embed apps
var AppsFS embed.FS

func (s *testService) configure() test.Test {
	return s.test("configure", func(t *testing.T) {
		s.Service = &portald.Service[target.Portal_]{}
		s.Config = s.config
		s.Config.Node.Log.Level = 100
		s.ExtraTokens = []string{"portal"}
		s.AppSources = []target.Source{source.Embed(AppsFS)}
		err := s.Configure()
		test.AssertErr(t, err)
		//s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald} // Faster testing
		s.Astrald = &debug.Astrald{NodeRoot: s.Config.Astrald} // Debugging astrald

		f := bind.CoreFactory{Repository: *s.Tokens()}
		s.Resolve = target.Any[target.Runnable](goja.Runner(f.NewBackendFunc()).Try)
	})
}

func (s *testService) start() test.Test {
	return s.test("start", func(t *testing.T) {
		err := s.Start(s.ctx)
		test.AssertErr(t, err)
	})
}

func (s *testService) nodeAlias() test.Test {
	return s.test("get node alias", func(t *testing.T) {
		alias, err := s.Apphost.NodeAlias()
		test.AssertErr(t, err)
		assert.NotZero(t, alias)
		s.alias = alias
	})
}

func (s *testService) createUser() test.Test {
	return s.test("create user", func(t *testing.T) {
		err := s.CreateUser("test_user")
		test.AssertErr(t, err)
	})
}

func (s *testService) userClaim(s2 *testService) test.Test {
	return s.test("claim", func(t *testing.T) {
		err := s.Claim(s2.Apphost.HostID.String())
		test.AssertErr(t, err)
	})
}

func (s *testService) addEndpoint(s2 *testService) test.Test {
	return s.test("add endpoint "+s2.name, func(t *testing.T) {
		id := s2.Apphost.HostID.String()
		port := s2.Config.TCP.ListenPort
		endpoint := fmt.Sprintf("tcp:127.0.0.1:%d", port)
		err := nodes.Client(&s.Apphost).AddEndpoint(id, endpoint)
		test.AssertErr(t, err)
	})
}

func (s *testService) installApps() test.Test {
	return s.test("install apps", func(t *testing.T) {
		l := source.Embed(apps.Builds)
		ctx := context.Background()
		for _, r := range s.Installer().Dispatcher().List(l) {
			err := r.Run(ctx)
			test.AssertErr(t, err)
			s.apps = append(s.apps, r)
		}
	})
}

func (s *testService) uninstallApp() test.Test {
	return s.test("uninstall app", func(t *testing.T) {
		if len(s.apps) == 0 {
			t.FailNow()
		}
		p := s.apps[0].Manifest().Package
		err := s.Installer().Uninstall(p)
		test.AssertErr(t, err)
	})
}

func (s *testService) publishAppBundles() test.Test {
	return s.test("publish app bundles", func(t *testing.T) {
		b := source.Embed(apps.Builds)
		s.published = map[object.ID]bundle.Release{}
		for _, o := range exec.ResolveBundle.List(b) {
			id, r, err := s.Publisher().Publish(o)
			test.AssertErr(t, err)
			s.published[*id] = *r
		}
	})
}

func (s *testService) awaitPublishedBundles() test.Test {
	return s.test("await bundles", func(t *testing.T) {
		time.Sleep(2000 * time.Millisecond)
		for id := range s.published {
			s.testAwaitDescribe(t, id) // await fetching describe
			s.testReadObject(t, id)    // test read object
		}
	})
}

func (s *testService) testAwaitDescribe(t *testing.T, id object.ID) {
	t.Run(s.name+" await describe", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		args := objects.DescribeArgs{ID: id}
		limit := 10
		for {
			obj, err := c.Describe(args)
			if obj != nil {
				plog.Println(id, obj)
				break
			}
			test.AssertErr(t, err)
			if limit > 0 {
				limit--
			} else {
				t.FailNow()
			}
			time.Sleep(200 * time.Millisecond)
		}
	})
}

func (s *testService) testReadObject(t *testing.T, id object.ID) {
	t.Run(s.name+" read object", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		rc, err := c.Read(objects.ReadArgs{ID: id})
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		test.AssertErr(t, err)

		plog.Println(id, buf.String())
	})
}

func (s *testService) reconnectAsUser() test.Test {
	return s.test("reconnect user", func(t *testing.T) {
		s.Apphost.AuthToken = s.UserInfo.AccessToken
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) searchObjects(query string) test.Test {
	return s.test("search objects", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		search, err := c.Scan(objects.ScanArgs{
			Type: query,
			Zone: astral.AllZones,
		})
		test.AssertErr(t, err)

		count := 0
		for result := range search {
			count++
			plog.Println(result)
		}
		assert.Greater(t, count, 0)
	})
}

func (s *testService) fetchReleases() test.Test {
	return s.test("fetch app release", func(t *testing.T) {
		for id, release := range s.published {
			c := objects.Client(s.Apphost.Rpc())
			r := &bundle.Release{}
			err := c.Fetch(objects.ReadArgs{ID: id}, r)
			test.AssertErr(t, err)
			assert.Equal(t, *release.BundleID, *r.BundleID)
			assert.Equal(t, *release.ManifestID, *r.ManifestID)
			assert.Equal(t, release.Release, r.Release)
			assert.Equal(t, release.Target, r.Target)
		}
	})
}

func (s *testService) fetchAppBundleExecs() test.Test {
	return s.test("fetch app bundle exes", func(t *testing.T) {
		for _, r := range s.published {
			c := objects.Client(s.Apphost.Rpc())
			o := &bundle.Object[target.Exec]{}
			o.Resolve = exec.ResolveBundle
			err := c.Fetch(objects.ReadArgs{ID: *r.BundleID}, o)
			test.AssertErr(t, err)
		}
	})
}

func (s *testService) openApp(pkg string) test.Test {
	return s.test("open  "+pkg, func(t *testing.T) {
		o := portald2.OpenOpt{}
		err := s.Open().Run(s.ctx, o, pkg)
		test.AssertErr(t, err)
	})
}

func (s *testService) setupToken(pkg string) test.Test {
	return s.test("setup token  "+pkg, func(t *testing.T) {
		_, err := s.Tokens().Resolve(pkg)
		test.AssertErr(t, err)
	})
}
