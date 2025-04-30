package astrald

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/core/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
	"time"
)

func TestInitializer_Start(t *testing.T) {
	plog.Verbosity = plog.Debug
	dir := test.Dir(t)
	test.Clean(dir)
	for i := range 2 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			r := Initializer{}
			r.AgentAlias = "portald"
			r.NodeRoot = dir
			r.TokensDir = dir
			r.Apphost = &apphost.Adapter{}
			r.Runner = &exec.Astrald{NodeRoot: dir}
			r.Config.Node.Log.Level = 100
			r.Config.Apphost.Listen = []string{"tcp:127.0.0.1:8635"}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(func() {
				cancel()
				time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
			})
			err := r.Start(ctx)
			if err != nil {
				plog.New().Println(err)
				t.Error(err)
			}
		})
	}
}
