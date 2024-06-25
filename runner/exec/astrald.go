package exec

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/apphost"
	exec2 "github.com/cryptopunkscc/portal/pkg/exec"
	"os"
	"os/exec"
	"time"
)

// Astral starts astral daemon process in a given [context.Context] and waits until available.
func Astral(ctx context.Context) (err error) {

	// check if astrald already running
	if err = apphost.Check(); err == nil {
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
	return exec2.Retry(ctx, 10*time.Second, func(int, int, time.Duration) (err error) {
		return apphost.Init()
	})
}
