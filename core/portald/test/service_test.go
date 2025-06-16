package portald

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/nodes"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/portal"
	portald2 "github.com/cryptopunkscc/portal/api/portald"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/api/user"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/astrald/debug"
	"github.com/cryptopunkscc/portal/core/bind"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/runner/goja"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"runtime"
	"strings"
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
	published map[astral.ObjectID]bundle.Release
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
		s.Resolve = target.Any[target.Runnable](
			target.Skip("node_modules"),
			goja.Runner(f.NewBackendFunc()).Try,
		)
	})
}

func (s *testService) start() test.Test {
	return s.test(func(t *testing.T) {
		err := s.Start(s.ctx)
		test.AssertErr(t, err)
	})
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
	})
}

func (s *testService) userInfo() test.Test {
	return s.test(func(t *testing.T) {
		info, err := s.UserInfo()
		test.AssertErr(t, err)
		assert.NotZero(t, info)
		plog.Println(*info)
	})
}

func (s *testService) userClaim(s2 *testService) test.Test {
	return s.test(func(t *testing.T) {
		err := s.Claim(s2.Apphost.HostID.String())
		test.AssertErr(t, err)
	})
}

func (s *testService) addEndpoint(s2 *testService) test.Test {
	return s.arg(s2.name).test(func(t *testing.T) {
		id := s2.Apphost.HostID.String()
		port := s2.Config.TCP.ListenPort
		endpoint := fmt.Sprintf("tcp:127.0.0.1:%d", port)
		err := nodes.Op(&s.Apphost).AddEndpoint(id, endpoint)
		test.AssertErr(t, err)
	})
}

func (s *testService) buildApps() test.Test {
	return s.test(func(t *testing.T) {
		err := apps.Build("pack")
		if err.Error() != "npm is required but not installed" {
			test.AssertErr(t, err)
		}
	})
}

func (s *testService) installApps() test.Test {
	return s.test(func(t *testing.T) {
		l := source.Embed(apps.Builds)
		ctx := context.Background()
		for _, r := range s.Installer().Dispatcher().List(l) {
			err := r.Run(ctx)
			test.AssertErr(t, err)
			s.apps = append(s.apps, r)
		}
		assert.NotEmpty(t, s.apps)
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
		s.published = map[astral.ObjectID]bundle.Release{}

		l, err := s.PublishAppsFS(b)
		test.AssertErr(t, err)
		count := 0
		for _, app := range l {
			s.published[*app.ReleaseID] = app.Release
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

func (s *testService) awaitObject(id astral.ObjectID) test.Test {
	return s.test(func(t *testing.T) {
		c := objects.Op(s.Apphost.Rpc())
		args := objects.ReadArgs{ID: id}
		limit := 10
		for {
			rc, err := c.Read(args)
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
		c := objects.Op(s.Apphost.Rpc())
		rc, err := c.Read(objects.ReadArgs{ID: id})
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		test.AssertErr(t, err)

		plog.Println(id, buf.String())
	})
}

func (s *testService) reconnectAsUser() test.Test {
	return s.test(func(t *testing.T) {
		s.Apphost.AuthToken = s.UserCreated.AccessToken
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) reconnectAsUser2(s2 *testService) test.Test {
	return s.test(func(t *testing.T) {
		s.Apphost.AuthToken = s2.UserCreated.AccessToken
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) reconnectAs(alias string) test.Test {
	return s.arg(alias).test(func(t *testing.T) {
		pt, err := s.Tokens().Resolve(alias)
		test.AssertErr(t, err)
		s.Apphost.AuthToken = pt.Token.String()
		err = s.Apphost.Reconnect()
		test.AssertErr(t, err)
	})
}

func (s *testService) scanObjects(typ string, s2 ...*testService) test.Test {
	o := s2
	if len(o) == 0 {
		o = append(o, s)
	}
	return s.arg(typ).test(func(t *testing.T) {
		c := objects.Op(s.Apphost.Rpc(), o[0].Apphost.HostID.String())
		search, err := c.Scan(objects.ScanArgs{
			Type: typ,
			Zone: astral.ZoneAll,
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
		c := objects.Op(s.Apphost.Rpc(), o[0].Apphost.HostID.String())
		search, err := c.Search(objects.SearchArgs{
			Query: query,
			Zone:  astral.ZoneAll,
		})
		test.AssertErr(t, err)

		count := 0
		for result := range search {
			count++
			plog.Println(result)
		}
		assert.NotZero(t, count)
	})
}

func (s *testService) fetchReleases() test.Test {
	return s.test(func(t *testing.T) {
		for id, release := range s.published {
			c := objects.Op(s.Apphost.Rpc())
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

func (s *testService) signAppContract(pkg string) test.Test {
	return s.arg(pkg).test(func(t *testing.T) {
		id, err := s.Apphost.Resolve(pkg)
		test.AssertErr(t, err)
		contract, err := apphost.Op(&s.Apphost).SignAppContract(id)
		test.AssertErr(t, err)
		plog.Println(contract)
	})
}

func (s *testService) fetchAppBundleExecs() test.Test {
	return s.test(func(t *testing.T) {
		for _, r := range s.published {
			c := objects.Op(s.Apphost.Rpc())
			o := &bundle.Object[any]{}
			o.Resolve = target.Any[target.AppBundle[any]](bundle.Resolve_.Try)
			err := c.Fetch(objects.ReadArgs{ID: *r.BundleID}, o)
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
	})
}

func (s *testService) listSiblings() test.Test {
	return s.test(func(t *testing.T) {
		c := user.Op(&s.Apphost)
		siblings, err := c.Siblings()
		test.AssertErr(t, err)
		count := 0
		for sibling := range siblings {
			count++
			plog.Println(sibling.String())
		}
		assert.Equal(t, 1, count)
	})
}
