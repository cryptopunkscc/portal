package exec

import (
	"context"
	"github.com/cryptopunkscc/astrald/lib/apphost"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"os"
	"os/exec"
)

func WithReadWriter(ctx context.Context, rw io.ReadWriter) context.Context {
	return context.WithValue(ctx, stdKey, &Std{In: rw, Out: rw, Err: rw})
}

func RunCmd(ctx context.Context, token string, path string, args ...string) (err error) {
	log := plog.Get(ctx).Scope("exec.RunCmd").Set(&ctx)
	log.Printf("Command run: %s, %v", path, args)
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = append(os.Environ(), apphost.AuthTokenEnv+"="+token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err = setStd(cmd, ctx); err != nil {
		return
	}
	err = cmd.Run()
	log.Printf("Command finished: %s, %v, %v", path, args, err)
	return
}

func setStd(cmd *exec.Cmd, ctx context.Context) (err error) {
	log := plog.Get(ctx)
	std, ok := ctx.Value(stdKey).(*Std)
	log.Printf("redirecting std %v ", ok)
	if !ok {
		return
	}

	cmd.Stdin = nil
	cmd.Stdout = std.Out
	cmd.Stderr = std.Err
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	go func() {
		_, _ = io.Copy(stdIn, std.In)
		_ = stdIn.Close()
	}()
	return
}

const stdKey = "exec.RunCmd.Std"

type Std struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
