package portald

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/portal"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/core/portald/debug"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testService struct {
	name   string
	alias  string
	config portal.Config
	*Service[target.Portal_]
}

func (s *testService) setupDir(t *testing.T) {
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
		if err := s.Apphost.Nodes().AddEndpoint(id, endpoint); err != nil {
			plog.Println(err)
			t.FailNow()
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
