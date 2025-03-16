package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/runtime/apphost"
	"os"
	"os/exec"
	"time"
)

// Astral starts astral daemon process in a given [context.Context] and waits until available.
func Astral(ctx context.Context) (err error) {
	defer plog.TraceErr(&err)
	log := plog.Get(ctx)
	// check if astrald already running
	if err = apphost.Default.Connect(); err == nil {
		return
	}

	// start astrald process
	cmd := exec.CommandContext(ctx, "astrald")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err = cmd.Start(); err != nil {
		return
	}

	// await for apphost
	return flow.Retry(ctx, 10*time.Second, func(i int, n int, d time.Duration) (err error) {
		log.Printf("%d/%d attempt %v: retry after %v", i+1, n, err, d)
		return apphost.Default.Connect()
	})
}
