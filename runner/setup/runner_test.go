package setup

import (
	"context"
	"github.com/cryptopunkscc/portal/core/token"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/astrald"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestRunner_Setup(t *testing.T) {
	plog.Verbosity = plog.Debug
	testDir := test.Dir(t)
	r := Runner{
		NodeRoot: testDir,
		Tokens:   token.Repository{Dir: testDir},
		Runner: astrald.Runner{
			Astrald: &exec.Astrald{NodeRoot: testDir},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := r.Setup(ctx); err != nil {
		plog.New().Println(err)
		t.Fatal(err)
	}
}
