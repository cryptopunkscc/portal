package astrald

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
)

func TestRunner_Start(t *testing.T) {
	plog.Verbosity = plog.Debug
	testDir := test.Dir(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	r := Runner{Astrald: &exec.Astrald{NodeRoot: testDir}}
	err := r.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
