package portald

import (
	"context"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
	"time"
)

func TestService_Start(t *testing.T) {
	plog.Verbosity = plog.Debug
	nodeDir := mem.NewVar(test.Dir(t, ".test"))
	appsDir := mem.NewVar(test.Dir(t, ".test", "apps"))
	tokensDir := mem.NewVar(test.Dir(t, ".test", "tokens"))
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
	})

	s := Service[target.Portal_]{}
	s.NodeDir = nodeDir
	s.TokensDir = tokensDir
	s.AppsDir = appsDir
	s.Astrald = &exec.Astrald{NodeRoot: nodeDir}
	s.CreateTokens = []string{"portal"}

	err := s.Start(ctx)
	if err != nil {
		plog.New().Println(err)
		t.Error(err)
	}
}
