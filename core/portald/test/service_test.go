package portald

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/object"
	"github.com/cryptopunkscc/portal/api/nodes"
	"github.com/cryptopunkscc/portal/api/objects"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/apps"
	"github.com/cryptopunkscc/portal/core/astrald/debug"
	"github.com/cryptopunkscc/portal/core/portald"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/cryptopunkscc/portal/target/bundle"
	"github.com/cryptopunkscc/portal/target/exec"
	"github.com/cryptopunkscc/portal/target/source"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testService struct {
	name   string
	alias  string
	config portal.Config
	*portald.Service[target.Portal_]
	apps      []target.Portal_
	published map[object.ID]bundle.Release
}

func (s *testService) cleanDir(t *testing.T) {
	s.config.Dir = test.CleanDir(t, s.name)
}

func (s *testService) configure(t *testing.T) {
	t.Run(s.name+" configure", func(t *testing.T) {
		s.Service = &portald.Service[target.Portal_]{}
		s.Config = s.config
		s.Config.Node.Log.Level = 100
		s.ExtraTokens = []string{"portal"}
		err := s.Configure()
		test.AssertErr(t, err)
		//s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald} // Faster testing
		s.Astrald = &debug.Astrald{NodeRoot: s.Config.Astrald} // Debugging astrald
	})
}

func (s *testService) testNodeStart(t *testing.T, ctx context.Context) {
	t.Run(s.name+" start", func(t *testing.T) {
		err := s.Start(ctx)
		test.AssertErr(t, err)
	})
}

func (s *testService) testNodeAlias(t *testing.T) {
	t.Run(s.name+" get node alias", func(t *testing.T) {
		alias, err := s.Apphost.NodeAlias()
		test.AssertErr(t, err)
		assert.NotZero(t, alias)
		s.alias = alias
	})
}

func (s *testService) testCreateUser(t *testing.T) {
	t.Run(s.name+" create user", func(t *testing.T) {
		err := s.CreateUser("test_user")
		test.AssertErr(t, err)
	})
}

func (s *testService) testUserClaim(t *testing.T, s2 *testService) {
	t.Run(s.name+" claim", func(t *testing.T) {
		err := s.Claim(s2.Apphost.HostID.String())
		test.AssertErr(t, err)
	})
}

func (s *testService) testAddEndpoint(t *testing.T, s2 *testService) {
	t.Run(s.name+" add endpoint", func(t *testing.T) {
		id := s2.Apphost.HostID.String()
		port := s2.Config.TCP.ListenPort
		endpoint := fmt.Sprintf("tcp:127.0.0.1:%d", port)
		err := nodes.Client(&s.Apphost).AddEndpoint(id, endpoint)
		test.AssertErr(t, err)
	})
}

func (s *testService) testWriteObject(t *testing.T, obj astral.Object) (id *object.ID) {
	t.Run(s.name+" write object", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		_, err := astral.WriteCanonical(buf, obj)
		test.AssertErr(t, err)

		id, err = astral.ResolveObjectID(obj)
		test.AssertErr(t, err)

		path := filepath.Join(s.Config.Astrald, "data", id.String())
		err = os.WriteFile(path, buf.Bytes(), 0644)
	})
	return
}

func (s *testService) testReconnectAsUser(t *testing.T) {
	t.Run(s.name+" reconnect user", func(t *testing.T) {
		s.Apphost.AuthToken = s.UserInfo.AccessToken
		err := s.Apphost.Reconnect()
		test.AssertErr(t, err)
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

func (s *testService) testShowObject(t *testing.T, id object.ID) {
	t.Run(s.name+" show object", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		objStr, err := c.Show(id)
		test.AssertErr(t, err)
		plog.Println(id, objStr)
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

func (s *testService) testSearchObjects(t *testing.T, query string) {
	t.Run(s.name+" search objects", func(t *testing.T) {
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

func (s *testService) testInstallApps(t *testing.T) {
	t.Run(s.name+" install apps", func(t *testing.T) {
		l := source.Embed(apps.Builds)
		ctx := context.Background()
		for _, r := range s.Installer().Dispatcher().List(l) {
			err := r.Run(ctx)
			test.AssertErr(t, err)
			s.apps = append(s.apps, r)
		}
	})
}

func (s *testService) testUninstallApp(t *testing.T) {
	t.Run(s.name+" uninstall app", func(t *testing.T) {
		if len(s.apps) == 0 {
			t.FailNow()
		}
		p := s.apps[0].Manifest().Package
		err := s.Installer().Uninstall(p)
		test.AssertErr(t, err)
	})
}

func (s *testService) testPublishAppBundle(t *testing.T) {
	t.Run(s.name+" publish app bundles", func(t *testing.T) {
		b := source.Embed(apps.Builds)
		s.published = map[object.ID]bundle.Release{}
		for _, o := range exec.ResolveBundle.List(b) {
			id, r, err := s.Publisher().Publish(o)
			test.AssertErr(t, err)
			s.published[*id] = *r
		}
	})
}

func (s *testService) awaitPublishedObjects(t *testing.T) {
	t.Run(s.name+" await published objects", func(t *testing.T) {
		for id := range s.published {
			s.testAwaitDescribe(t, id) // await fetching describe
			s.testReadObject(t, id)    // test read object
		}
	})
}

func (s *testService) testFetchReleases(t *testing.T) {
	t.Run(s.name+" fetch app release", func(t *testing.T) {
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

func (s *testService) testFetchAppBundleExecs(t *testing.T) {
	t.Run(s.name+" fetch app bundle", func(t *testing.T) {
		for id := range s.published {
			c := objects.Client(s.Apphost.Rpc())
			o := &bundle.Object[target.Exec]{}
			o.Resolve = exec.ResolveBundle
			err := c.Fetch(objects.ReadArgs{ID: id}, o)
			test.AssertErr(t, err)
		}
	})
}

func testServiceContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
	})
	return ctx
}
