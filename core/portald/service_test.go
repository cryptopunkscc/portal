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
	"github.com/cryptopunkscc/portal/core/portald/debug"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
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
	*Service[target.Portal_]
}

func (s *testService) cleanDir(t *testing.T) {
	s.config.Dir = test.CleanDir(t, s.name)
}

func (s *testService) configure(t *testing.T) {
	t.Run(s.name+" configure", func(t *testing.T) {
		s.Service = &Service[target.Portal_]{}
		s.Config = s.config
		s.Config.Node.Log.Level = 100
		s.ExtraTokens = []string{"portal"}
		if err := s.Configure(); err != nil {
			plog.Println(err)
			t.FailNow()
		}
		//s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald} // Faster testing
		s.Astrald = &debug.Astrald{NodeRoot: s.Config.Astrald} // Debugging astrald
	})
}

func (s *testService) testNodeStart(t *testing.T, ctx context.Context) {
	t.Run(s.name+" start", func(t *testing.T) {
		if err := s.Start(ctx); err != nil {
			plog.Println(err)
			t.FailNow()
		}
	})
}

func (s *testService) testNodeAlias(t *testing.T) {
	t.Run(s.name+" get node alias", func(t *testing.T) {
		if alias, err := s.Apphost.NodeAlias(); err != nil {
			plog.Println(err)
			t.FailNow()
		} else {
			assert.NotZero(t, alias)
			s.alias = alias
		}
	})
}

func (s *testService) testCreateUser(t *testing.T) {
	t.Run(s.name+" create user", func(t *testing.T) {
		if err := s.CreateUser("test_user"); err != nil {
			plog.Println(err)
			t.FailNow()
		}
	})
}

func (s *testService) testUserClaim(t *testing.T, s2 *testService) {
	t.Run(s.name+" claim", func(t *testing.T) {
		if err := s.Claim(s2.Apphost.HostID.String()); err != nil {
			plog.Println(err)
			t.FailNow()
		}
	})
}

func (s *testService) testAddEndpoint(t *testing.T, s2 *testService) {
	t.Run(s.name+" add endpoint", func(t *testing.T) {
		id := s2.Apphost.HostID.String()
		port := s2.Config.TCP.ListenPort
		endpoint := fmt.Sprintf("tcp:127.0.0.1:%d", port)
		if err := nodes.Client(&s.Apphost).AddEndpoint(id, endpoint); err != nil {
			//if err := s.Apphost.Client().Nodes().AddEndpoint(id, endpoint); err != nil {
			plog.Println(err)
			t.FailNow()
		}
	})
}

func (s *testService) testWriteObject(t *testing.T, obj astral.Object) (id object.ID) {
	t.Run(s.name+" write object", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		_, err := astral.WriteCanonical(buf, obj)
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
		id, err = astral.ResolveObjectID(obj)
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
		path := filepath.Join(s.Config.Astrald, "data", id.String())
		err = os.WriteFile(path, buf.Bytes(), 0644)
	})
	return
}

func (s *testService) testReconnectAsUser(t *testing.T) {
	t.Run(s.name+" reconnect user", func(t *testing.T) {
		s.Apphost.AuthToken = s.UserInfo.AccessToken
		err := s.Apphost.Reconnect()
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
	})
}

func (s *testService) testAwaitDescribe(t *testing.T, id object.ID) {
	t.Run(s.name+" await describe", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		args := objects.DescribeArgs{ID: id}
		for {
			time.Sleep(200 * time.Millisecond)
			if obj, err := c.Describe(args); obj != nil {
				plog.Println(id, obj)
				break
			} else if err != nil {
				plog.Println(err)
				t.FailNow()
			}
		}
	})
}

func (s *testService) testShowObject(t *testing.T, id object.ID) {
	t.Run(s.name+" show object", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		objStr, err := c.Show(id)
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
		plog.Println(id, objStr)
	})
}

func (s *testService) testReadObject(t *testing.T, id object.ID) {
	t.Run(s.name+" read object", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		rc, err := c.Read(objects.ReadArgs{ID: id})
		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(rc)
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
		plog.Println(id, buf.String())
	})
}

func (s *testService) testSearchObjects(t *testing.T, query string) {
	t.Run(s.name+" search objects", func(t *testing.T) {
		c := objects.Client(s.Apphost.Rpc())
		search, err := c.Search(objects.SearchArgs{
			Query: query,
			Zone:  astral.AllZones,
		})
		if err != nil {
			plog.Println(err)
			t.FailNow()
		}
		count := 0
		for result := range search {
			count++
			plog.Println(result)
		}
		assert.Equal(t, 1, count)
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
