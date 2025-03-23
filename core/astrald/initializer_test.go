package astrald

import (
	"context"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/mem"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
	"time"
)

func TestRunner_Start(t *testing.T) {
	plog.Verbosity = plog.Debug
	testDir := mem.NewVar(test.Dir(t))
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
		time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
	})

	r := Initializer{}
	r.NodeRoot = testDir
	r.TokensDir = testDir
	r.Apphost = &apphost.Adapter{}
	r.Runner = &exec.Astrald{NodeRoot: testDir}

	err := r.Start(ctx)
	if err != nil {
		plog.New().Println(err)
		t.Error(err)
	}
}
