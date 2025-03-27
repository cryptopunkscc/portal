package portald

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runner/exec"
	"github.com/cryptopunkscc/portal/test"
	"testing"
	"time"
)

func TestService_Start(t *testing.T) {
	plog.Verbosity = plog.Debug
	dir := test.Dir(t)
	test.Clean(dir)
	for i := range 2 {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			s := Service[target.Portal_]{}
			s.Config.Dir = dir
			s.Config.Node.Log.Level = 100
			s.Config.Apphost.Listen = []string{"tcp:127.0.0.1:8635"}
			if err := s.Config.Build(); err != nil {
				plog.New().Println(err)
				t.Error(err)
			}
			s.Astrald = &exec.Astrald{NodeRoot: s.Config.Astrald}
			s.ExtraTokens = []string{"portal"}

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(func() {
				cancel()
				time.Sleep(10 * time.Millisecond) // give a time to kill astrald process
			})
			if err := s.Start(ctx); err != nil {
				plog.New().Println(err)
				t.Error(err)
			}
		})
	}
}
