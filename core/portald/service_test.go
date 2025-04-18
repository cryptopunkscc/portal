package portald

import (
	"context"
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
	config portal.Config
	*Service[target.Portal_]
}

func (s *testService) setupDir(t *testing.T) {
	s.config.Dir = test.CleanDir(t, s.name)
}

func (s *testService) configure() {
	s.Service = &Service[target.Portal_]{}
	s.Config = s.config
	s.Config.Node.Log.Level = 100
	s.ExtraTokens = []string{"portal"}
	if err := s.Configure(); err != nil {
		plog.P().Println(err)
	}
	//s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald} // Faster testing
	s.Astrald = &debug.Astrald{NodeRoot: s.Config.Astrald} // Debugging astrald
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
